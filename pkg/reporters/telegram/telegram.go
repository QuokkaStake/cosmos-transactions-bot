package telegram

import (
	"errors"
	"fmt"
	"html"
	"main/pkg/constants"
	"main/pkg/data_fetcher"
	"main/pkg/metrics"
	"main/pkg/templates"
	"main/pkg/types"
	"strings"
	"time"

	"gopkg.in/telebot.v3/middleware"

	"main/pkg/alias_manager"
	"main/pkg/config"
	configTypes "main/pkg/config/types"
	nodesManager "main/pkg/nodes_manager"
	"main/pkg/utils"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v3"
)

type Reporter struct {
	ReporterName string

	Token  string
	Chat   int64
	Admins []int64

	TelegramBot      *tele.Bot
	Logger           zerolog.Logger
	NodesManager     *nodesManager.NodesManager
	Config           *config.AppConfig
	AliasManager     *alias_manager.AliasManager
	MetricsManager   *metrics.Manager
	DataFetcher      *data_fetcher.DataFetcher
	TemplatesManager templates.Manager

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
		ReporterName:     reporterConfig.Name,
		Token:            reporterConfig.TelegramConfig.Token,
		Chat:             reporterConfig.TelegramConfig.Chat,
		Admins:           reporterConfig.TelegramConfig.Admins,
		Config:           config,
		Logger:           logger.With().Str("component", "telegram_reporter").Logger(),
		TemplatesManager: templates.NewTelegramTemplateManager(logger, config.Timezone),
		NodesManager:     nodesManager,
		AliasManager:     aliasManager,
		MetricsManager:   metricsManager,
		DataFetcher:      dataFetcher,
		Version:          version,
		StopChannel:      make(chan bool),
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

func (reporter *Reporter) SerializeReport(r types.Report) (string, error) {
	reportableType := r.Reportable.Type()
	return reporter.TemplatesManager.Render(reportableType, r)
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

func (reporter *Reporter) Stop() {
	reporter.StopChannel <- true
}
