package telegram

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"main/pkg/types/amount"
	"time"

	"github.com/dustin/go-humanize"
	"gopkg.in/telebot.v3/middleware"

	"main/pkg/alias_manager"
	"main/pkg/config"
	configTypes "main/pkg/config/types"
	nodesManager "main/pkg/nodes_manager"
	"main/pkg/types"
	"main/pkg/utils"
	"main/templates"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v3"
)

type TelegramReporter struct {
	Token  string
	Chat   int64
	Admins []int64

	TelegramBot  *tele.Bot
	Logger       zerolog.Logger
	Templates    map[string]*template.Template
	NodesManager *nodesManager.NodesManager
	Config       *config.AppConfig
	AliasManager *alias_manager.AliasManager
}

const (
	MaxMessageSize = 4096
)

func NewTelegramReporter(
	config *config.AppConfig,
	logger *zerolog.Logger,
	nodesManager *nodesManager.NodesManager,
	aliasManager *alias_manager.AliasManager,
) *TelegramReporter {
	return &TelegramReporter{
		Token:        config.TelegramConfig.Token,
		Chat:         config.TelegramConfig.Chat,
		Admins:       config.TelegramConfig.Admins,
		Config:       config,
		Logger:       logger.With().Str("component", "telegram_reporter").Logger(),
		Templates:    make(map[string]*template.Template, 0),
		NodesManager: nodesManager,
		AliasManager: aliasManager,
	}
}

func (reporter *TelegramReporter) Init() {
	if reporter.Token == "" || reporter.Chat == 0 {
		reporter.Logger.Debug().Msg("Telegram credentials not set, not creating Telegram reporter.")
		return
	}

	bot, err := tele.NewBot(tele.Settings{
		Token:  reporter.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		reporter.Logger.Warn().Err(err).Msg("Could not create Telegram bot")
		return
	}

	if len(reporter.Admins) > 0 {
		reporter.Logger.Debug().Msg("Using admins whitelist")
		bot.Use(middleware.Whitelist(reporter.Admins...))
	}

	bot.Handle("/help", reporter.HandleHelp)
	bot.Handle("/start", reporter.HandleHelp)
	bot.Handle("/status", reporter.HandleListNodesStatus)
	bot.Handle("/config", reporter.HandleGetConfig)
	bot.Handle("/alias", reporter.HandleSetAlias)
	bot.Handle("/aliases", reporter.HandleGetAliases)

	reporter.TelegramBot = bot
	go reporter.TelegramBot.Start()
}

func (reporter TelegramReporter) Enabled() bool {
	return reporter.Token != "" && reporter.Chat != 0
}

func (reporter TelegramReporter) GetTemplate(name string) (*template.Template, error) {
	if cachedTemplate, ok := reporter.Templates[name]; ok {
		reporter.Logger.Trace().Str("type", name).Msg("Using cached template")
		return cachedTemplate, nil
	}

	reporter.Logger.Trace().Str("type", name).Msg("Loading template")

	filename := fmt.Sprintf("%s.html", utils.RemoveFirstSlash(name))

	t, err := template.New(filename).Funcs(template.FuncMap{
		"SerializeLink":    reporter.SerializeLink,
		"SerializeAmount":  reporter.SerializeAmount,
		"SerializeDate":    reporter.SerializeDate,
		"SerializeMessage": reporter.SerializeMessage,
	}).ParseFS(templates.TemplatesFs, "telegram/"+filename)
	if err != nil {
		return nil, err
	}

	reporter.Templates[name] = t

	return t, nil
}

func (reporter *TelegramReporter) Render(templateName string, data interface{}) (string, error) {
	template, err := reporter.GetTemplate(templateName)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", templateName).Msg("Error loading template")
		return "", err
	}

	var buffer bytes.Buffer
	err = template.Execute(&buffer, data)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", templateName).Msg("Error rendering template")
		return "", err
	}

	return buffer.String(), err
}

func (reporter *TelegramReporter) SerializeReport(r types.Report) (string, error) {
	reportableType := r.Reportable.Type()
	return reporter.Render(reportableType, r)
}

func (reporter *TelegramReporter) SerializeMessage(msg types.Message) template.HTML {
	msgType := msg.Type()

	reporterTemplate, err := reporter.GetTemplate(msgType)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", msgType).Msg("Error loading template")
		return template.HTML(fmt.Sprintf("Error loading template: <code>%s</code>", html.EscapeString(err.Error())))
	}

	var buffer bytes.Buffer
	err = reporterTemplate.Execute(&buffer, msg)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", msgType).Msg("Error rendering template")
		return template.HTML(fmt.Sprintf("Error rendering template: <code>%s</code>", html.EscapeString(err.Error())))
	}

	return template.HTML(buffer.String())
}

func (reporter TelegramReporter) Send(report types.Report) error {
	reportString, err := reporter.SerializeReport(report)
	if err != nil {
		reporter.Logger.Error().
			Err(err).
			Msg("Could not serialize Telegram message to report, trying to send fallback message")

		if err := reporter.BotSend("Error serializing report, check logs for more info."); err != nil {
			reporter.Logger.Err(err).Msg("Could not send Telegram fallback message")
			return err
		}

		return nil
	}

	reporter.Logger.Trace().Str("report", reportString).Msg("Sending a report")

	if err := reporter.BotSend(reportString); err != nil {
		reporter.Logger.Err(err).Msg("Could not send Telegram message")
		return err
	}
	return nil
}

func (reporter TelegramReporter) Name() string {
	return "telegram-reporter"
}

func (reporter *TelegramReporter) BotSend(msg string) error {
	messages := utils.SplitStringIntoChunks(msg, MaxMessageSize)

	for _, message := range messages {
		if _, err := reporter.TelegramBot.Send(
			&tele.User{
				ID: reporter.Chat,
			},
			message,
			tele.ModeHTML,
			tele.NoPreview,
		); err != nil {
			reporter.Logger.Error().Err(err).Msg("Could not send Telegram message")
			return err
		}
	}
	return nil
}

func (reporter *TelegramReporter) BotReply(c tele.Context, msg string) error {
	messages := utils.SplitStringIntoChunks(msg, MaxMessageSize)

	for _, message := range messages {
		if err := c.Reply(message, tele.ModeHTML); err != nil {
			reporter.Logger.Error().Err(err).Msg("Could not send Telegram message")
			return err
		}
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

func (reporter *TelegramReporter) SerializeAmount(amount amount.Amount) template.HTML {
	if amount.PriceUSD == nil {
		return template.HTML(fmt.Sprintf(
			"%s %s",
			utils.StripTrailingDigits(humanize.BigCommaf(amount.Value), 6),
			amount.Denom,
		))
	}

	return template.HTML(fmt.Sprintf(
		"%s %s ($%s)",
		utils.StripTrailingDigits(humanize.BigCommaf(amount.Value), 6),
		amount.Denom,
		utils.StripTrailingDigits(humanize.BigCommaf(amount.PriceUSD), 3),
	))
}

func (reporter *TelegramReporter) SerializeDate(date time.Time) template.HTML {
	return template.HTML(date.Format(time.RFC822))
}
