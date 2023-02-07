package pkg

import (
	"main/pkg/types"
	"os"
	"os/signal"
	"syscall"

	"main/pkg/alias_manager"
	"main/pkg/config"
	"main/pkg/data_fetcher"
	"main/pkg/filterer"
	"main/pkg/logger"
	nodesManager "main/pkg/nodes_manager"
	"main/pkg/reporters"
	"main/pkg/reporters/telegram"

	"github.com/rs/zerolog"
)

type App struct {
	Logger       zerolog.Logger
	NodesManager *nodesManager.NodesManager
	Reporters    []reporters.Reporter
	DataFetchers map[string]*data_fetcher.DataFetcher
	Filterers    map[string]*filterer.Filterer
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

	filterers := make(map[string]*filterer.Filterer, len(config.Chains))
	for _, chain := range config.Chains {
		filterers[chain.Name] = filterer.NewFilterer(log, chain)
	}

	return &App{
		Logger:       log.With().Str("component", "app").Logger(),
		Reporters:    reporters,
		NodesManager: nodesManager,
		DataFetchers: dataFetchers,
		Filterers:    filterers,
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
		case rawReport := <-a.NodesManager.Channel:
			chainFilterer, _ := a.Filterers[rawReport.Chain.Name]
			fetcher, _ := a.DataFetchers[rawReport.Chain.Name]

			reportableFiltered := chainFilterer.Filter(rawReport.Reportable)
			if reportableFiltered == nil {
				a.Logger.Debug().
					Str("node", rawReport.Node).
					Str("chain", rawReport.Chain.Name).
					Str("hash", rawReport.Reportable.GetHash()).
					Msg("Got report")
				continue
			}

			report := types.Report{
				Node:       rawReport.Node,
				Chain:      rawReport.Chain,
				Reportable: reportableFiltered,
			}

			a.Logger.Info().
				Str("node", report.Node).
				Str("chain", report.Chain.Name).
				Str("hash", report.Reportable.GetHash()).
				Msg("Got report")

			rawReport.Reportable.GetAdditionalData(*fetcher)

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
