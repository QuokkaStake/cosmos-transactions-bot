package pkg

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"main/pkg/alias_manager"
	"main/pkg/config"
	"main/pkg/data_fetcher"
	"main/pkg/logger"
	nodesManager "main/pkg/nodes_manager"
	"main/pkg/reporters"
	"main/pkg/reporters/telegram"
)

type App struct {
	Logger       zerolog.Logger
	NodesManager *nodesManager.NodesManager
	Reporters    []reporters.Reporter
	DataFetchers map[string]*data_fetcher.DataFetcher
}

func NewApp(config *config.AppConfig) *App {
	log := logger.GetLogger(config.LogConfig)
	aliasManager := alias_manager.NewAliasManager(log, config)
	aliasManager.Load()

	nodesManager := nodesManager.NewNodesManager(log, config)

	reporters := []reporters.Reporter{
		telegram.NewTelegramReporter(config, log, nodesManager, aliasManager),
	}

	dataFetchers := make(map[string]*data_fetcher.DataFetcher, len(config.Chains))
	for _, chain := range config.Chains {
		dataFetchers[chain.Name] = data_fetcher.NewDataFetcher(log, chain, aliasManager)
	}

	return &App{
		Logger:       log.With().Str("component", "app").Logger(),
		Reporters:    reporters,
		NodesManager: nodesManager,
		DataFetchers: dataFetchers,
	}
}

func (a *App) Start() {
	for _, reporter := range a.Reporters {
		reporter.Init()
		if reporter.Enabled() {
			a.Logger.Info().Str("name", reporter.Name()).Msg("Init reporter")
		}
	}

	a.NodesManager.Listen()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case report := <-a.NodesManager.Channel:
			a.Logger.Info().
				Str("node", report.Node).
				Str("chain", report.Chain.Name).
				Str("hash", report.Reportable.GetHash()).
				Msg("Got report")

			fetcher, _ := a.DataFetchers[report.Chain.Name]
			report.Reportable.GetAdditionalData(*fetcher)

			for _, reporter := range a.Reporters {
				if err := reporter.Send(report); err != nil {
					a.Logger.Error().
						Err(err).
						Msg("Error sending report")
				}
			}
		case <-quit:
			a.NodesManager.Stop()
			os.Exit(0)
		}
	}
}
