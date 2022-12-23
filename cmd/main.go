package main

import (
	"github.com/spf13/cobra"
	"main/pkg"
	"main/pkg/config"
	"main/pkg/logger"
)

func Execute(configPath string) {
	config, err := config.GetConfig(configPath)
	if err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not load config")
	}

	if err = config.Validate(); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Provided config is invalid!")
	}

	app := pkg.NewApp(config)
	app.Start()
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
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not set flag as required")
	}

	if err := rootCmd.Execute(); err != nil {
		logger.GetDefaultLogger().Fatal().Err(err).Msg("Could not start application")
	}
}
