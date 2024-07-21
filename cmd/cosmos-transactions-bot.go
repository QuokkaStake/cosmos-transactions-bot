package main

import (
	"main/pkg"
	configPkg "main/pkg/config"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"

	"github.com/spf13/cobra"
)

var (
	version = "unknown"
)

func ExecuteMain(configPath string) {
	filesystem := &fs.OsFS{}
	app := pkg.NewApp(filesystem, configPath, version)
	app.Start()
}

func ExecuteValidateConfig(configPath string) {
	filesystem := &fs.OsFS{}

	config, err := configPkg.GetConfig(configPath, filesystem)
	if err != nil {
		loggerPkg.GetDefaultLogger().Panic().Err(err).Msg("Could not load config!")
	}

	if warnings := config.DisplayWarnings(); len(warnings) > 0 {
		for _, warning := range warnings {
			warning.Log(loggerPkg.GetDefaultLogger())
		}
	}

	loggerPkg.GetDefaultLogger().Info().Msg("Provided config is valid.")
}

func main() {
	var ConfigPath string

	rootCmd := &cobra.Command{
		Use:     "cosmos-transactions-bot --config [config path]",
		Long:    "Get notified on new transactions on different cosmos-sdk chains.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			ExecuteMain(ConfigPath)
		},
	}

	validateConfigCmd := &cobra.Command{
		Use:     "validate-config --config [config path]",
		Long:    "Validate config.",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			ExecuteValidateConfig(ConfigPath)
		},
	}

	rootCmd.PersistentFlags().StringVar(&ConfigPath, "config", "", "Config file path")
	_ = rootCmd.MarkPersistentFlagRequired("config")

	validateConfigCmd.PersistentFlags().StringVar(&ConfigPath, "config", "", "Config file path")
	_ = validateConfigCmd.MarkPersistentFlagRequired("config")

	rootCmd.AddCommand(validateConfigCmd)

	if err := rootCmd.Execute(); err != nil {
		loggerPkg.GetDefaultLogger().Panic().Err(err).Msg("Could not start application")
	}
}
