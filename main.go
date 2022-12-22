package main

import (
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func Execute(configPath string) {
	config, err := GetConfig(configPath)
	if err != nil {
		GetDefaultLogger().Fatal().Err(err).Msg("Could not load config")
	}

	if err = config.Validate(); err != nil {
		GetDefaultLogger().Fatal().Err(err).Msg("Provided config is invalid!")
	}

	log := GetLogger(config.LogConfig)

	nodesManager := NewNodesManager(log, config)
	nodesManager.Listen()

	reporters := []Reporter{
		NewTelegramReporter(config.TelegramConfig, log),
	}

	for _, reporter := range reporters {
		reporter.Init()
		if reporter.Enabled() {
			log.Info().Str("name", reporter.Name()).Msg("Init reporter")
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case report := <-nodesManager.Channel:
			log.Info().
				Str("node", report.Node).
				Str("chain", report.Chain.Name).
				Str("hash", report.Reportable.GetHash()).
				Msg("Got report")

			for _, reporter := range reporters {
				if err := reporter.Send(report); err != nil {
					log.Error().
						Err(err).
						Msg("Error sending report")
				}
			}
		case <-quit:
			nodesManager.Stop()
			os.Exit(0)
		}
	}
}

func main() {
	var ConfigPath string

	rootCmd := &cobra.Command{
		Use:  "cosmos-transactions-bot",
		Long: "Get notified on new transactions on different cosmos-sdk chains.",
		Run: func(cmd *cobra.Command, args []string) {
			Execute(ConfigPath)
		},
	}

	rootCmd.PersistentFlags().StringVar(&ConfigPath, "config", "", "Config file path")
	if err := rootCmd.MarkPersistentFlagRequired("config"); err != nil {
		GetDefaultLogger().Fatal().Err(err).Msg("Could not set flag as required")
	}

	if err := rootCmd.Execute(); err != nil {
		GetDefaultLogger().Fatal().Err(err).Msg("Could not start application")
	}
}
