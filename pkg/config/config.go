package config

import (
	"os"

	tomlConfig "main/pkg/config/toml_config"
	"main/pkg/config/types"
	"main/pkg/utils"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
	"github.com/rs/zerolog"
)

type Chains []*types.Chain

func (c Chains) FindByName(name string) *types.Chain {
	for _, chain := range c {
		if chain.Name == name {
			return chain
		}
	}

	return nil
}

type AppConfig struct {
	Path           string
	TelegramConfig TelegramConfig
	LogConfig      LogConfig
	Chains         Chains
}

type TelegramConfig struct {
	Chat   int64
	Token  string
	Admins []int64
}

type LogConfig struct {
	LogLevel   string
	JSONOutput bool
}

func GetConfig(path string) (*AppConfig, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configString := string(configBytes)

	configStruct := &tomlConfig.TomlConfig{}
	if _, err = toml.Decode(configString, configStruct); err != nil {
		return nil, err
	}

	if err = defaults.Set(configStruct); err != nil {
		return nil, err
	}

	if err := configStruct.Validate(); err != nil {
		return nil, err
	}

	return FromTomlConfig(configStruct, path), nil
}

func FromTomlConfig(c *tomlConfig.TomlConfig, path string) *AppConfig {
	return &AppConfig{
		Path: path,
		TelegramConfig: TelegramConfig{
			Chat:   c.TelegramConfig.Chat,
			Token:  c.TelegramConfig.Token,
			Admins: c.TelegramConfig.Admins,
		},
		LogConfig: LogConfig{
			LogLevel:   c.LogConfig.LogLevel,
			JSONOutput: c.LogConfig.JSONOutput,
		},
		Chains: utils.Map(c.Chains, func(c *tomlConfig.Chain) *types.Chain {
			return c.ToAppConfigChain()
		}),
	}
}

func (c *AppConfig) DisplayWarnings(log *zerolog.Logger) {
	for _, chain := range c.Chains {
		chain.DisplayWarnings(log)
	}
}
