package main

import (
	"main/pkg"
	configPkg "main/pkg/config"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "unknown"
)

type OsFS struct {
}

func (fs *OsFS) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (fs *OsFS) Create(path string) (fs.File, error) {
	return os.Create(path)
}

func Execute(configPath string) {
	filesystem := &OsFS{}
	config, err := configPkg.GetConfig(configPath, filesystem)
	if err != nil {
		loggerPkg.GetDefaultLogger().Fatal().Err(err).Msg("Could not load config")
	}
	warnings := config.DisplayWarnings()

	for _, warning := range warnings {
		warning.Log(loggerPkg.GetDefaultLogger())
	}

	app := pkg.NewApp(config, version, filesystem)
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
