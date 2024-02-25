package pkg

import (
	configTypes "main/pkg/config/types"
	"os"
	"os/signal"
	"syscall"

	"main/pkg/alias_manager"
	"main/pkg/config"
	"main/pkg/data_fetcher"
	filtererPkg "main/pkg/filterer"
	loggerPkg "main/pkg/logger"
	metricsPkg "main/pkg/metrics"
	nodesManagerPkg "main/pkg/nodes_manager"
	reportersPkg "main/pkg/reporters"

	"github.com/rs/zerolog"
)

type App struct {
	Logger         zerolog.Logger
	Config         *config.AppConfig
	Chains         []*configTypes.Chain
	NodesManager   *nodesManagerPkg.NodesManager
	Reporters      reportersPkg.Reporters
	DataFetcher    *data_fetcher.DataFetcher
	Filterer       *filtererPkg.Filterer
	MetricsManager *metricsPkg.Manager

	Version string
}

func NewApp(config *config.AppConfig, version string) *App {
	logger := loggerPkg.GetLogger(config.LogConfig)
	aliasManager := alias_manager.NewAliasManager(logger, config)
	aliasManager.Load()

	metricsManager := metricsPkg.NewManager(logger, config.Metrics)

	nodesManager := nodesManagerPkg.NewNodesManager(logger, config, metricsManager)

	reporters := make([]reportersPkg.Reporter, len(config.Reporters))
	for index, reporterConfig := range config.Reporters {
		reporters[index] = reportersPkg.GetReporter(
			reporterConfig,
			config,
			logger,
			nodesManager,
			aliasManager,
			version,
		)
	}

	dataFetcher := data_fetcher.NewDataFetcher(
		logger,
		config,
		aliasManager,
		metricsManager,
	)

	filterer := filtererPkg.NewFilterer(logger, config, metricsManager)

	return &App{
		Logger:         logger.With().Str("component", "app").Logger(),
		Config:         config,
		Chains:         config.Chains,
		Reporters:      reporters,
		NodesManager:   nodesManager,
		DataFetcher:    dataFetcher,
		Filterer:       filterer,
		MetricsManager: metricsManager,
		Version:        version,
	}
}

func (a *App) Start() {
	go a.MetricsManager.Start()
	a.MetricsManager.LogAppVersion(a.Version)
	a.MetricsManager.SetAllDefaultMetrics(a.Config)

	for _, reporter := range a.Reporters {
		reporter.Init()
		a.MetricsManager.LogReporterEnabled(reporter.Name(), reporter.Type())
		a.Logger.Info().
			Str("name", reporter.Name()).
			Str("type", reporter.Type()).
			Msg("Init reporter")
	}

	a.NodesManager.Listen()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case rawReport := <-a.NodesManager.Channel:
			reportablesForReporters := a.Filterer.GetReportableForReporters(rawReport)

			if len(reportablesForReporters) == 0 {
				a.Logger.Debug().
					Str("node", rawReport.Node).
					Str("chain", rawReport.Chain.Name).
					Str("hash", rawReport.Reportable.GetHash()).
					Msg("Got report which is nowhere to send")
				continue
			}

			for reporterName, report := range reportablesForReporters {
				a.Logger.Info().
					Str("node", report.Node).
					Str("chain", report.Chain.Name).
					Str("reporter", reporterName).
					Str("hash", report.Reportable.GetHash()).
					Msg("Got report")

				rawReport.Reportable.GetAdditionalData(a.DataFetcher)

				reporter := a.Reporters.FindByName(reporterName)

				if err := reporter.Send(report); err != nil {
					a.Logger.Error().
						Err(err).
						Msg("Error sending report")
					a.MetricsManager.LogReport(report, reporterName, false)
				} else {
					a.MetricsManager.LogReport(report, reporterName, true)
				}
			}
		case <-quit:
			a.NodesManager.Stop()
			os.Exit(0)
		}
	}
}
