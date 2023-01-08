package config

import (
	"bytes"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
	"github.com/rs/zerolog"
	tomlConfig "main/pkg/config/toml_config"
	"main/pkg/config/types"
	"main/pkg/utils"
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
	AliasesPath    string
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
		Path:        path,
		AliasesPath: c.AliasesPath,
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

func (c *AppConfig) ToTomlConfig() *tomlConfig.TomlConfig {
	return &tomlConfig.TomlConfig{
		AliasesPath: c.AliasesPath,
		TelegramConfig: tomlConfig.TelegramConfig{
			Chat:   c.TelegramConfig.Chat,
			Token:  c.TelegramConfig.Token,
			Admins: c.TelegramConfig.Admins,
		},
		LogConfig: tomlConfig.LogConfig{
			LogLevel:   c.LogConfig.LogLevel,
			JSONOutput: c.LogConfig.JSONOutput,
		},
		Chains: utils.Map(c.Chains, tomlConfig.FromAppConfigChain),
	}
}

func (c *AppConfig) DisplayWarnings(log *zerolog.Logger) {
	for _, chain := range c.Chains {
		chain.DisplayWarnings(log)
	}
}

func (c *AppConfig) Save() error {
	configStruct := c.ToTomlConfig()

	f, err := os.Create(c.Path)
	if err != nil {
		return err
	}
	if err := toml.NewEncoder(f).Encode(configStruct); err != nil {
		return err
	}
	return f.Close()
}

func (c *AppConfig) GetConfigAsString() (string, error) {
	configStruct := c.ToTomlConfig()

	buffer := new(bytes.Buffer)

	if err := toml.NewEncoder(buffer).Encode(configStruct); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
