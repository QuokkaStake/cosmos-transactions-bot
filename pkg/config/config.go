package config

import (
	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
	tomlConfig "main/pkg/config/toml_config"
	"main/pkg/config/types"
	"main/pkg/utils"
	"os"
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
	TelegramChat  int64
	TelegramToken string
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
			TelegramChat:  c.TelegramConfig.TelegramChat,
			TelegramToken: c.TelegramConfig.TelegramToken,
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
