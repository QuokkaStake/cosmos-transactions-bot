package telegram

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"main/pkg/config"
	configTypes "main/pkg/config/types"
	nodesManager "main/pkg/nodes_manager"
	"main/pkg/types"
	"main/pkg/utils"
	"strings"
	"time"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v3"
)

type TelegramReporter struct {
	TelegramToken string
	TelegramChat  int64

	TelegramBot  *tele.Bot
	Logger       zerolog.Logger
	Templates    map[string]*template.Template
	NodesManager *nodesManager.NodesManager
}

const (
	MaxMessageSize = 4096
)

//go:embed templates/*
var templatesFs embed.FS

type TelegramSerializedReport struct {
	Report types.Report
	Msgs   []template.HTML
}

func NewTelegramReporter(
	config config.TelegramConfig,
	logger *zerolog.Logger,
	nodesManager *nodesManager.NodesManager,
) *TelegramReporter {
	return &TelegramReporter{
		TelegramToken: config.TelegramToken,
		TelegramChat:  config.TelegramChat,
		Logger:        logger.With().Str("component", "telegram_reporter").Logger(),
		Templates:     make(map[string]*template.Template, 0),
		NodesManager:  nodesManager,
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

	bot.Handle("/status", reporter.HandleListNodesStatus)

	reporter.TelegramBot = bot
	go reporter.TelegramBot.Start()
}

func (reporter TelegramReporter) Enabled() bool {
	return reporter.TelegramToken != "" && reporter.TelegramChat != 0
}

func (reporter TelegramReporter) GetTemplate(name string) (*template.Template, error) {
	if template, ok := reporter.Templates[name]; ok {
		reporter.Logger.Trace().Str("type", name).Msg("Using cached template")
		return template, nil
	}

	reporter.Logger.Trace().Str("type", name).Msg("Loading template")

	filename := fmt.Sprintf("%s.html", utils.RemoveFirstSlash(name))

	t, err := template.New(filename).Funcs(template.FuncMap{
		"SerializeLink":    reporter.SerializeLink,
		"SerializeAmount":  reporter.SerializeAmount,
		"SerializeDate":    reporter.SerializeDate,
		"SerializeMessage": reporter.SerializeMessage,
	}).ParseFS(templatesFs, "templates/"+filename)
	if err != nil {
		return nil, err
	}

	reporter.Templates[name] = t

	return t, nil
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

func (reporter *TelegramReporter) SerializeMessage(msg types.Message) template.HTML {
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

func (reporter TelegramReporter) Send(report types.Report) error {
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
		if sb.Len()+len(line) >= MaxMessageSize {
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

func (reporter *TelegramReporter) SerializeLink(link configTypes.Link) template.HTML {
	value := link.Title
	if value == "" {
		value = link.Value
	}

	if link.Href != "" {
		return template.HTML(fmt.Sprintf("<a href='%s'>%s</a>", link.Href, value))
	}

	return template.HTML(value)
}

func (reporter *TelegramReporter) SerializeAmount(amount types.Amount) template.HTML {
	if amount.PriceUSD != 0 {
		return template.HTML(fmt.Sprintf(
			"%.6f%s ($%.2f)",
			amount.Value,
			amount.Denom,
			amount.PriceUSD,
		))
	}

	return template.HTML(fmt.Sprintf(
		"%.6f%s",
		amount.Value,
		amount.Denom,
	))
}

func (reporter *TelegramReporter) SerializeDate(date time.Time) template.HTML {
	return template.HTML(date.Format(time.RFC822))
}
