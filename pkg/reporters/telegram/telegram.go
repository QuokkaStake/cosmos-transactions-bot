package telegram

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"html/template"
	"main/pkg/constants"
	"main/pkg/data_fetcher"
	"main/pkg/metrics"
	"main/pkg/types"
	"main/pkg/types/amount"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"gopkg.in/telebot.v3/middleware"

	"main/pkg/alias_manager"
	"main/pkg/config"
	configTypes "main/pkg/config/types"
	nodesManager "main/pkg/nodes_manager"
	"main/pkg/utils"
	"main/templates"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v3"
)

type Reporter struct {
	ReporterName string

	Token  string
	Chat   int64
	Admins []int64

	TelegramBot    *tele.Bot
	Logger         zerolog.Logger
	Templates      map[string]*template.Template
	NodesManager   *nodesManager.NodesManager
	Config         *config.AppConfig
	AliasManager   *alias_manager.AliasManager
	MetricsManager *metrics.Manager
	DataFetcher    *data_fetcher.DataFetcher

	Version     string
	StopChannel chan bool
}

const (
	MaxMessageSize = 4096
)

func NewReporter(
	reporterConfig *configTypes.Reporter,
	config *config.AppConfig,
	logger *zerolog.Logger,
	nodesManager *nodesManager.NodesManager,
	aliasManager *alias_manager.AliasManager,
	metricsManager *metrics.Manager,
	dataFetcher *data_fetcher.DataFetcher,
	version string,
) *Reporter {
	return &Reporter{
		ReporterName:   reporterConfig.Name,
		Token:          reporterConfig.TelegramConfig.Token,
		Chat:           reporterConfig.TelegramConfig.Chat,
		Admins:         reporterConfig.TelegramConfig.Admins,
		Config:         config,
		Logger:         logger.With().Str("component", "telegram_reporter").Logger(),
		Templates:      make(map[string]*template.Template),
		NodesManager:   nodesManager,
		AliasManager:   aliasManager,
		MetricsManager: metricsManager,
		DataFetcher:    dataFetcher,
		Version:        version,
		StopChannel:    make(chan bool),
	}
}

func (reporter *Reporter) Init() error {
	if reporter.Token == "" || reporter.Chat == 0 {
		reporter.Logger.Debug().Msg("Telegram credentials not set, not creating Telegram reporter.")
		return nil
	}

	bot, err := tele.NewBot(tele.Settings{
		Token:  reporter.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		reporter.Logger.Warn().Err(err).Msg("Could not create Telegram bot")
		return err
	}

	if len(reporter.Admins) > 0 {
		reporter.Logger.Debug().Msg("Using admins whitelist")
		bot.Use(middleware.Whitelist(reporter.Admins...))
	}

	reporter.AddCommand("/help", bot, reporter.GetHelpCommand())
	reporter.AddCommand("/start", bot, reporter.GetHelpCommand())
	reporter.AddCommand("/status", bot, reporter.GetListNodesCommand())
	reporter.AddCommand("/alias", bot, reporter.GetSetAliasCommand())
	reporter.AddCommand("/aliases", bot, reporter.GetGetAliasesCommand())

	reporter.TelegramBot = bot

	return nil
}

func (reporter *Reporter) Start() {
	if reporter.TelegramBot == nil {
		return
	}

	go reporter.TelegramBot.Start()

	<-reporter.StopChannel
	reporter.Logger.Info().Msg("Shutting down...")
	reporter.TelegramBot.Stop()
}

func (reporter *Reporter) AddCommand(query string, bot *tele.Bot, command Command) {
	bot.Handle(query, reporter.Handler(command))
}

func (reporter *Reporter) Handler(command Command) tele.HandlerFunc {
	return func(c tele.Context) error {
		reporter.Logger.Info().
			Str("sender", c.Sender().Username).
			Str("text", c.Text()).
			Str("command", command.Name).
			Msg("Got query")

		reporter.MetricsManager.LogReporterQuery(reporter.Name(), command.Query)

		args := strings.Split(c.Text(), " ")

		if len(args)-1 < command.MinArgs {
			if err := reporter.BotReply(c, html.EscapeString(fmt.Sprintf(command.Usage, args[0]))); err != nil {
				return err
			}

			return errors.New("invalid invocation")
		}

		result, err := command.Execute(c)
		if err != nil {
			reporter.Logger.Error().
				Err(err).
				Str("command", command.Name).
				Msg("Error processing command")
			if result != "" {
				return reporter.BotReply(c, result)
			} else {
				return reporter.BotReply(c, "Internal error!")
			}
		}

		return reporter.BotReply(c, result)
	}
}

func (reporter *Reporter) GetTemplate(name string) (*template.Template, error) {
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

func (reporter *Reporter) Render(templateName string, data interface{}) (string, error) {
	reportTemplate, err := reporter.GetTemplate(templateName)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", templateName).Msg("Error loading template")
		return "", err
	}

	var buffer bytes.Buffer
	err = reportTemplate.Execute(&buffer, data)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", templateName).Msg("Error rendering template")
		return "", err
	}

	return buffer.String(), err
}

func (reporter *Reporter) SerializeReport(r types.Report) (string, error) {
	reportableType := r.Reportable.Type()
	return reporter.Render(reportableType, r)
}

func (reporter *Reporter) SerializeMessage(msg types.Message) template.HTML {
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

func (reporter *Reporter) Send(report types.Report) error {
	reportString, err := reporter.SerializeReport(report)
	if err != nil {
		reporter.Logger.Error().
			Err(err).
			Msg("Could not serialize Telegram message to report, trying to send fallback message")

		if sendErr := reporter.BotSend("Error serializing report, check logs for more info."); sendErr != nil {
			reporter.Logger.Err(sendErr).Msg("Could not send Telegram fallback message")
			return sendErr
		}

		return nil
	}

	reporter.Logger.Trace().Str("report", reportString).Msg("Sending a report")

	if sendErr := reporter.BotSend(reportString); sendErr != nil {
		reporter.Logger.Err(sendErr).Msg("Could not send Telegram message")
		return sendErr
	}
	return nil
}

func (reporter *Reporter) Name() string {
	return reporter.ReporterName
}

func (reporter *Reporter) Type() string {
	return constants.ReporterTypeTelegram
}

func (reporter *Reporter) BotSend(msg string) error {
	messages := utils.SplitStringIntoChunks(msg, MaxMessageSize)

	for _, message := range messages {
		if _, err := reporter.TelegramBot.Send(
			&tele.User{ID: reporter.Chat},
			strings.TrimSpace(message),
			tele.ModeHTML,
			tele.NoPreview,
		); err != nil {
			reporter.Logger.Error().Err(err).Msg("Could not send Telegram message")
			return err
		}
	}
	return nil
}

func (reporter *Reporter) BotReply(c tele.Context, msg string) error {
	messages := utils.SplitStringIntoChunks(msg, MaxMessageSize)

	for _, message := range messages {
		if err := c.Reply(strings.TrimSpace(message), tele.ModeHTML, tele.NoPreview); err != nil {
			reporter.Logger.Error().Err(err).Msg("Could not send Telegram message")
			return err
		}
	}
	return nil
}

func (reporter *Reporter) SerializeLink(link *configTypes.Link) template.HTML {
	value := link.Title
	if value == "" {
		value = link.Value
	}

	if link.Href != "" {
		return template.HTML(fmt.Sprintf("<a href='%s'>%s</a>", link.Href, value))
	}

	return template.HTML(value)
}

func (reporter *Reporter) SerializeAmount(amount amount.Amount) template.HTML {
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

func (reporter *Reporter) SerializeDate(date time.Time) template.HTML {
	return template.HTML(date.In(reporter.Config.Timezone).Format(time.RFC822))
}

func (reporter *Reporter) Stop() {
	reporter.StopChannel <- true
}
