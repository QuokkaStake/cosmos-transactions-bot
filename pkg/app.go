package pkg

import (
	configTypes "main/pkg/config/types"
	fsPkg "main/pkg/fs"
	"main/pkg/types"
	"os"
	"os/signal"
	"syscall"

	"main/pkg/alias_manager"
	configPkg "main/pkg/config"
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
	Config         *configPkg.AppConfig
	Chains         []*configTypes.Chain
	NodesManager   *nodesManagerPkg.NodesManager
	Reporters      reportersPkg.Reporters
	DataFetcher    *data_fetcher.DataFetcher
	Filterer       *filtererPkg.Filterer
	MetricsManager *metricsPkg.Manager
	QuitChannel    chan os.Signal

	Version string
}

func NewApp(filesystem fsPkg.FS, configPath string, version string) *App {
	config, err := configPkg.GetConfig(configPath, filesystem)
	if err != nil {
		loggerPkg.GetDefaultLogger().Panic().Err(err).Msg("Could not load config")
	}
	warnings := config.DisplayWarnings()

	for _, warning := range warnings {
		warning.Log(loggerPkg.GetDefaultLogger())
	}

	logger := loggerPkg.GetLogger(config.LogConfig)
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	aliasManager.Load()

	metricsManager := metricsPkg.NewManager(logger, config.Metrics)
	nodesManager := nodesManagerPkg.NewNodesManager(logger, config, metricsManager)
	dataFetcher := data_fetcher.NewDataFetcher(
		logger,
		config,
		aliasManager,
		metricsManager,
	)

	reporters := make([]reportersPkg.Reporter, len(config.Reporters))
	for index, reporterConfig := range config.Reporters {
		reporters[index] = reportersPkg.GetReporter(
			reporterConfig,
			config,
			logger,
			nodesManager,
			aliasManager,
			metricsManager,
			dataFetcher,
			version,
		)
	}

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
		QuitChannel:    make(chan os.Signal, 1),
	}
}

func (a *App) Start() {
	go a.MetricsManager.Start()
	a.MetricsManager.LogAppVersion(a.Version)
	a.MetricsManager.SetAllDefaultMetrics(a.Config)

	for _, reporter := range a.Reporters {
		if err := reporter.Init(); err != nil {
			continue
		}

		go reporter.Start()
		a.MetricsManager.LogReporterEnabled(reporter.Name(), reporter.Type())
		a.Logger.Info().
			Str("name", reporter.Name()).
			Str("type", reporter.Type()).
			Msg("Init reporter")
	}

	a.NodesManager.Listen()

	signal.Notify(a.QuitChannel, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case rawReport := <-a.NodesManager.Channel:
			a.ProcessReport(rawReport)
		case <-a.QuitChannel:
			a.NodesManager.Stop()
			a.MetricsManager.Stop()
			return
		}
	}
}

func (a *App) ProcessReport(rawReport types.Report) {
	reportablesForReporters := a.Filterer.GetReportableForReporters(rawReport)

	if len(reportablesForReporters) == 0 {
		a.Logger.Debug().
			Str("node", rawReport.Node).
			Str("chain", rawReport.Chain.Name).
			Str("hash", rawReport.Reportable.GetHash()).
			Msg("Got report which is nowhere to send")
	}

	for reporterName, report := range reportablesForReporters {
		a.Logger.Info().
			Str("node", report.Node).
			Str("chain", report.Chain.Name).
			Str("reporter", reporterName).
			Str("hash", report.Reportable.GetHash()).
			Msg("Got report")

		report.Reportable.GetAdditionalData(a.DataFetcher, report.Subscription.Name)

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
}
