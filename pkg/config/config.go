package config

import (
	"bytes"
	"main/pkg/fs"
	"time"

	"gopkg.in/guregu/null.v4"

	tomlConfig "main/pkg/config/toml_config"
	"main/pkg/config/types"
	"main/pkg/utils"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
)

type AppConfig struct {
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

type ReadFileFs interface {
	ReadFile(name string) ([]byte, error)
}

func GetConfig(path string, filesystem fs.FS) (*AppConfig, error) {
	configBytes, err := filesystem.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configString := string(configBytes)

	configStruct := &tomlConfig.TomlConfig{}
	if _, err = toml.Decode(configString, configStruct); err != nil {
		return nil, err
	}

	defaults.MustSet(configStruct)

	if err := configStruct.Validate(); err != nil {
		return nil, err
	}

	return FromTomlConfig(configStruct), nil
}

func FromTomlConfig(c *tomlConfig.TomlConfig) *AppConfig {
	timezone, _ := time.LoadLocation(c.Timezone)

	return &AppConfig{
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

func (c *AppConfig) DisplayWarnings() []types.DisplayWarning {
	var warnings []types.DisplayWarning

	for _, chain := range c.Chains {
		warnings = append(warnings, chain.DisplayWarnings()...)
	}

	reportersUsed := map[string]bool{}
	chainsUsed := map[string]bool{}

	for _, subscription := range c.Subscriptions {
		reportersUsed[subscription.Reporter] = true

		for _, chainSubscription := range subscription.ChainSubscriptions {
			chainsUsed[chainSubscription.Chain] = true
		}
	}

	for _, chain := range c.Chains {
		if _, ok := chainsUsed[chain.Name]; !ok {
			warnings = append(warnings, types.DisplayWarning{
				Keys: map[string]string{"chain": chain.Name},
				Text: "Chain is not used in any subscriptions",
			})
		}
	}

	for _, reporter := range c.Reporters {
		if _, ok := reportersUsed[reporter.Name]; !ok {
			warnings = append(warnings, types.DisplayWarning{
				Keys: map[string]string{"reporter": reporter.Name},
				Text: "Reporter is not used in any subscriptions",
			})
		}
	}

	return warnings
}

func (c *AppConfig) GetConfigAsString() string {
	configStruct := c.ToTomlConfig()
	buffer := new(bytes.Buffer)
	_ = toml.NewEncoder(buffer).Encode(configStruct)
	return buffer.String()
}
