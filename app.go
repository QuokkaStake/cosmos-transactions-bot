package main

import (
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	Logger       zerolog.Logger
	NodesManager *NodesManager
	Reporters    []Reporter
	DataFetchers map[string]*DataFetcher
}

func NewApp(config *Config) *App {
	logger := GetLogger(config.LogConfig)

	nodesManager := NewNodesManager(logger, config)

	reporters := []Reporter{
		NewTelegramReporter(config.TelegramConfig, logger),
	}

	dataFetchers := make(map[string]*DataFetcher, len(config.Chains))
	for _, chain := range config.Chains {
		dataFetchers[chain.Name] = NewDataFetcher(logger, chain)
	}

	return &App{
		Logger:       logger.With().Str("component", "telegram_reporter").Logger(),
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
