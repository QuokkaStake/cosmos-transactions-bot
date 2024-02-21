package config

import (
	"bytes"
	"os"
	"time"

	"gopkg.in/guregu/null.v4"

	tomlConfig "main/pkg/config/toml_config"
	"main/pkg/config/types"
	"main/pkg/utils"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
	"github.com/rs/zerolog"
)

type AppConfig struct {
	Path          string
	AliasesPath   string
	LogConfig     LogConfig
	Chains        types.Chains
	Subscriptions types.Subscriptions
	Reporters     types.Reporters
	Metrics       MetricsConfig
	Timezone      *time.Location
}

type LogConfig struct {
	LogLevel   string
	JSONOutput bool
}

type MetricsConfig struct {
	Enabled    bool
	ListenAddr string
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
	timezone, _ := time.LoadLocation(c.Timezone)

	return &AppConfig{
		Path:        path,
		AliasesPath: c.AliasesPath,
		LogConfig: LogConfig{
			LogLevel:   c.LogConfig.LogLevel,
			JSONOutput: c.LogConfig.JSONOutput.Bool,
		},
		Metrics: MetricsConfig{
			ListenAddr: c.MetricsConfig.ListenAddr,
			Enabled:    c.MetricsConfig.Enabled.Bool,
		},
		Chains: utils.Map(c.Chains, func(c *tomlConfig.Chain) *types.Chain {
			return c.ToAppConfigChain()
		}),
		Reporters: utils.Map(c.Reporters, func(r *tomlConfig.Reporter) *types.Reporter {
			return r.ToAppConfigReporter()
		}),
		Subscriptions: utils.Map(c.Subscriptions, func(s *tomlConfig.Subscription) *types.Subscription {
			return s.ToAppConfigSubscription()
		}),
		Timezone: timezone,
	}
}

func (c *AppConfig) ToTomlConfig() *tomlConfig.TomlConfig {
	return &tomlConfig.TomlConfig{
		AliasesPath: c.AliasesPath,
		LogConfig: tomlConfig.LogConfig{
			LogLevel:   c.LogConfig.LogLevel,
			JSONOutput: null.BoolFrom(c.LogConfig.JSONOutput),
		},
		MetricsConfig: tomlConfig.MetricsConfig{
			ListenAddr: c.Metrics.ListenAddr,
			Enabled:    null.BoolFrom(c.Metrics.Enabled),
		},
		Chains:        utils.Map(c.Chains, tomlConfig.FromAppConfigChain),
		Reporters:     utils.Map(c.Reporters, tomlConfig.FromAppConfigReporter),
		Subscriptions: utils.Map(c.Subscriptions, tomlConfig.FromAppConfigSubscription),
		Timezone:      c.Timezone.String(),
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
