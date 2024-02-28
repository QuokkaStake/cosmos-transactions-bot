package toml_config_test

import (
	tomlConfig "main/pkg/config/toml_config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTomlConfigNoChains(t *testing.T) {
	t.Parallel()

	config := tomlConfig.TomlConfig{}
	require.Error(t, config.Validate())
}

func TestTomlConfigInvalidTimezone(t *testing.T) {
	t.Parallel()

	config := tomlConfig.TomlConfig{
		Chains: tomlConfig.Chains{
			{},
		},
		Timezone: "invalid",
	}
	require.Error(t, config.Validate())
}

func TestTomlConfigInvalidChain(t *testing.T) {
	t.Parallel()

	config := tomlConfig.TomlConfig{
		Chains: tomlConfig.Chains{
			{},
		},
		Timezone: "Etc/UTC",
	}
	require.Error(t, config.Validate())
}

func TestTomlConfigInvalidReporter(t *testing.T) {
	t.Parallel()

	config := tomlConfig.TomlConfig{
		Chains: tomlConfig.Chains{
			{
				Name:            "chain",
				ChainID:         "chain-id",
				TendermintNodes: []string{"node"},
				APINodes:        []string{"node"},
				Queries:         []string{"event.key = 'value'"},
			},
		},
		Reporters: tomlConfig.Reporters{
			{},
		},
		Timezone: "Etc/UTC",
	}
	require.Error(t, config.Validate())
}

func TestTomlConfigInvalidSubscription(t *testing.T) {
	t.Parallel()

	config := tomlConfig.TomlConfig{
		Chains: tomlConfig.Chains{
			{
				Name:            "chain",
				ChainID:         "chain-id",
				TendermintNodes: []string{"node"},
				APINodes:        []string{"node"},
				Queries:         []string{"event.key = 'value'"},
			},
		},
		Reporters: tomlConfig.Reporters{
			{
				Name: "test",
				Type: "telegram",
				TelegramConfig: &tomlConfig.TelegramConfig{
					Chat:   1,
					Token:  "xxx:yyy",
					Admins: []int64{123},
				},
			},
		},
		Subscriptions: tomlConfig.Subscriptions{
			{},
		},
		Timezone: "Etc/UTC",
	}
	require.Error(t, config.Validate())
}

func TestTomlConfigChainSubscriptionChainNotFound(t *testing.T) {
	t.Parallel()

	config := tomlConfig.TomlConfig{
		Chains: tomlConfig.Chains{
			{
				Name:            "chain",
				ChainID:         "chain-id",
				TendermintNodes: []string{"node"},
				APINodes:        []string{"node"},
				Queries:         []string{"event.key = 'value'"},
			},
		},
		Reporters: tomlConfig.Reporters{
			{
				Name: "test",
				Type: "telegram",
				TelegramConfig: &tomlConfig.TelegramConfig{
					Chat:   1,
					Token:  "xxx:yyy",
					Admins: []int64{123},
				},
			},
		},
		Subscriptions: tomlConfig.Subscriptions{
			{
				Name:     "name",
				Reporter: "reporter",
				ChainSubscriptions: tomlConfig.ChainSubscriptions{
					{Chain: "nonexistent"},
				},
			},
		},
		Timezone: "Etc/UTC",
	}
	require.Error(t, config.Validate())
}

func TestTomlConfigSubscriptionReporterNotFound(t *testing.T) {
	t.Parallel()

	config := tomlConfig.TomlConfig{
		Chains: tomlConfig.Chains{
			{
				Name:            "chain",
				ChainID:         "chain-id",
				TendermintNodes: []string{"node"},
				APINodes:        []string{"node"},
				Queries:         []string{"event.key = 'value'"},
			},
		},
		Reporters: tomlConfig.Reporters{
			{
				Name: "test",
				Type: "telegram",
				TelegramConfig: &tomlConfig.TelegramConfig{
					Chat:   1,
					Token:  "xxx:yyy",
					Admins: []int64{123},
				},
			},
		},
		Subscriptions: tomlConfig.Subscriptions{
			{
				Name:     "name",
				Reporter: "nonexistent",
				ChainSubscriptions: tomlConfig.ChainSubscriptions{
					{Chain: "chain"},
				},
			},
		},
		Timezone: "Etc/UTC",
	}
	require.Error(t, config.Validate())
}

func TestTomlConfigValid(t *testing.T) {
	t.Parallel()

	config := tomlConfig.TomlConfig{
		Chains: tomlConfig.Chains{
			{
				Name:            "chain",
				ChainID:         "chain-id",
				TendermintNodes: []string{"node"},
				APINodes:        []string{"node"},
				Queries:         []string{"event.key = 'value'"},
			},
		},
		Reporters: tomlConfig.Reporters{
			{
				Name: "test",
				Type: "telegram",
				TelegramConfig: &tomlConfig.TelegramConfig{
					Chat:   1,
					Token:  "xxx:yyy",
					Admins: []int64{123},
				},
			},
		},
		Subscriptions: tomlConfig.Subscriptions{
			{
				Name:     "name",
				Reporter: "test",
				ChainSubscriptions: tomlConfig.ChainSubscriptions{
					{Chain: "chain"},
				},
			},
		},
		Timezone: "Etc/UTC",
	}
	require.NoError(t, config.Validate())
}
