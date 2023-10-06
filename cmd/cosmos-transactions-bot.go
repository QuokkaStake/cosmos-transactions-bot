package main

import (
	"main/pkg"
	configPkg "main/pkg/config"
	loggerPkg "main/pkg/logger"

	"github.com/spf13/cobra"
)

var (
	version = "unknown"
)

func Execute(configPath string) {
	config, err := configPkg.GetConfig(configPath)
	if err != nil {
		loggerPkg.GetDefaultLogger().Fatal().Err(err).Msg("Could not load config")
	}
	config.DisplayWarnings(loggerPkg.GetLogger(config.LogConfig))

	app := pkg.NewApp(config, version)
	app.Start()
}

func main() {
	var ConfigPath string

	rootCmd := &cobra.Command{
		Use:     "cosmos-transactions-bot",
		Long:    "Get notified on new transactions on different cosmos-sdk chains.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			Execute(ConfigPath)
		},
	}

	rootCmd.PersistentFlags().StringVar(&ConfigPath, "config", "", "Config file path")
	if err := rootCmd.MarkPersistentFlagRequired("config"); err != nil {
		loggerPkg.GetDefaultLogger().Fatal().Err(err).Msg("Could not set flag as required")
	}

	if err := rootCmd.Execute(); err != nil {
		loggerPkg.GetDefaultLogger().Fatal().Err(err).Msg("Could not start application")
	}
}
