package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v3"
)

type TelegramReporter struct {
	TelegramToken string
	TelegramChat  int64

	TelegramBot *tele.Bot
	Logger      zerolog.Logger
	Templates   map[string]*template.Template
}

const (
	MaxMessageSize = 4096
)

//go:embed templates/telegram*
var templatesFs embed.FS

type TelegramSerializedReport struct {
	Report Report
	Msgs   []template.HTML
}

func NewTelegramReporter(
	config TelegramConfig,
	logger *zerolog.Logger) *TelegramReporter {
	return &TelegramReporter{
		TelegramToken: config.TelegramToken,
		TelegramChat:  config.TelegramChat,
		Logger:        logger.With().Str("component", "telegram_reporter").Logger(),
		Templates:     make(map[string]*template.Template, 0),
	}
}

func (reporter *TelegramReporter) Init() {
	if reporter.TelegramToken == "" || reporter.TelegramChat == 0 {
		reporter.Logger.Debug().Msg("Telegram credentials not set, not creating Telegram reporter.")
		return
	}

	bot, err := tele.NewBot(tele.Settings{
		Token:  reporter.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		reporter.Logger.Warn().Err(err).Msg("Could not create Telegram bot")
		return
	}

	reporter.TelegramBot = bot
	go reporter.TelegramBot.Start()
}

func (reporter TelegramReporter) Enabled() bool {
	return reporter.TelegramToken != "" && reporter.TelegramChat != 0
}

func (reporter TelegramReporter) GetTemplate(t string) (*template.Template, error) {
	if template, ok := reporter.Templates[t]; ok {
		reporter.Logger.Trace().Str("type", string(t)).Msg("Using cached template")
		return template, nil
	}

	reporter.Logger.Trace().Str("type", string(t)).Msg("Loading template")

	filename := fmt.Sprintf("templates/telegram/%s.html", t)
	template, err := template.ParseFS(templatesFs, filename)
	if err != nil {
		return nil, err
	}

	reporter.Templates[t] = template

	return template, nil
}

func (reporter *TelegramReporter) SerializeReport(e TelegramSerializedReport) (string, error) {
	reportableType := e.Report.Reportable.Type()

	template, err := reporter.GetTemplate(reportableType)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", reportableType).Msg("Error loading template")
		return "", err
	}

	var buffer bytes.Buffer
	err = template.Execute(&buffer, e)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", reportableType).Msg("Error rendering template")
		return "", err
	}

	return buffer.String(), nil
}

func (reporter *TelegramReporter) SerializeMessage(msg Message) template.HTML {
	msgType := msg.Type()

	reporterTemplate, err := reporter.GetTemplate(msgType)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", msgType).Msg("Error loading template")
		return template.HTML(fmt.Sprintf("Error loading template: %s", err))
	}

	var buffer bytes.Buffer
	err = reporterTemplate.Execute(&buffer, msg)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", msgType).Msg("Error rendering template")
		return template.HTML(fmt.Sprintf("Error rendering template: %s", err))
	}

	return template.HTML(buffer.String())
}

func (reporter TelegramReporter) Send(report Report) error {
	msgsSerialized := make([]template.HTML, len(report.Reportable.GetMessages()))

	for index, msg := range report.Reportable.GetMessages() {
		msgsSerialized[index] = reporter.SerializeMessage(msg)
	}

	reportSerialized := TelegramSerializedReport{
		Report: report,
		Msgs:   msgsSerialized,
	}

	reportString, err := reporter.SerializeReport(reportSerialized)
	if err != nil {
		reporter.Logger.Error().Err(err).Msg("Could not serialize Telegram message to report")
		return err
	}

	_, err = reporter.TelegramBot.Send(
		&tele.User{
			ID: reporter.TelegramChat,
		},
		reportString,
		tele.ModeHTML,
		tele.NoPreview,
	)
	if err != nil {
		reporter.Logger.Err(err).Msg("Could not send Telegram message")
		return err
	}
	return nil
}

func (reporter TelegramReporter) Name() string {
	return "telegram-reporter"
}

func (reporter *TelegramReporter) BotReply(c tele.Context, msg string) error {
	msgsByNewline := strings.Split(msg, "\n")

	var sb strings.Builder

	for _, line := range msgsByNewline {
		if sb.Len()+len(line) > MaxMessageSize {
			if err := c.Reply(sb.String(), tele.ModeHTML); err != nil {
				reporter.Logger.Error().Err(err).Msg("Could not send Telegram message")
				return err
			}

			sb.Reset()
		}

		sb.WriteString(line + "\n")
	}

	if err := c.Reply(sb.String(), tele.ModeHTML); err != nil {
		reporter.Logger.Error().Err(err).Msg("Could not send Telegram message")
		return err
	}

	return nil
}
