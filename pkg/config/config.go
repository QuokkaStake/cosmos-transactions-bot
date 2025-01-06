package config

import (
	"main/pkg/fs"

	"gopkg.in/guregu/null.v4"

	"main/pkg/config/types"
	yamlConfig "main/pkg/config/yaml_config"
	"main/pkg/utils"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	AliasesPath   string
	LogConfig     LogConfig
	Chains        types.Chains
	Subscriptions types.Subscriptions
	Reporters     types.Reporters
	Metrics       MetricsConfig
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

	configStruct := &yamlConfig.YamlConfig{}
	if err = yaml.Unmarshal(configBytes, configStruct); err != nil {
		return nil, err
	}

	defaults.MustSet(configStruct)

	if err := configStruct.Validate(); err != nil {
		return nil, err
	}

	return FromYamlConfig(configStruct), nil
}

func FromYamlConfig(c *yamlConfig.YamlConfig) *AppConfig {
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
		Chains: utils.Map(c.Chains, func(c *yamlConfig.Chain) *types.Chain {
			return c.ToAppConfigChain()
		}),
		Reporters: utils.Map(c.Reporters, func(r *yamlConfig.Reporter) *types.Reporter {
			return r.ToAppConfigReporter()
		}),
		Subscriptions: utils.Map(c.Subscriptions, func(s *yamlConfig.Subscription) *types.Subscription {
			return s.ToAppConfigSubscription()
		}),
	}
}

func (c *AppConfig) ToYamlConfig() *yamlConfig.YamlConfig {
	return &yamlConfig.YamlConfig{
		AliasesPath: c.AliasesPath,
		LogConfig: yamlConfig.LogConfig{
			LogLevel:   c.LogConfig.LogLevel,
			JSONOutput: null.BoolFrom(c.LogConfig.JSONOutput),
		},
		MetricsConfig: yamlConfig.MetricsConfig{
			ListenAddr: c.Metrics.ListenAddr,
			Enabled:    null.BoolFrom(c.Metrics.Enabled),
		},
		Chains:        utils.Map(c.Chains, yamlConfig.FromAppConfigChain),
		Reporters:     utils.Map(c.Reporters, yamlConfig.FromAppConfigReporter),
		Subscriptions: utils.Map(c.Subscriptions, yamlConfig.FromAppConfigSubscription),
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
