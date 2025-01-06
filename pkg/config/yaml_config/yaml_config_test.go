package yaml_config_test

import (
	yamlConfig "main/pkg/config/yaml_config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestYamlConfigNoChains(t *testing.T) {
	t.Parallel()

	config := yamlConfig.YamlConfig{}
	require.Error(t, config.Validate())
}

func TestYamlConfigInvalidChain(t *testing.T) {
	t.Parallel()

	config := yamlConfig.YamlConfig{
		Chains: yamlConfig.Chains{
			{},
		},
	}
	require.Error(t, config.Validate())
}

func TestYamlConfigInvalidReporter(t *testing.T) {
	t.Parallel()

	config := yamlConfig.YamlConfig{
		Chains: yamlConfig.Chains{
			{
				Name:            "chain",
				ChainID:         "chain-id",
				TendermintNodes: []string{"node"},
				APINodes:        []string{"node"},
				Queries:         []string{"event.key = 'value'"},
			},
		},
		Reporters: yamlConfig.Reporters{
			{},
		},
	}
	require.Error(t, config.Validate())
}

func TestYamlConfigInvalidSubscription(t *testing.T) {
	t.Parallel()

	config := yamlConfig.YamlConfig{
		Chains: yamlConfig.Chains{
			{
				Name:            "chain",
				ChainID:         "chain-id",
				TendermintNodes: []string{"node"},
				APINodes:        []string{"node"},
				Queries:         []string{"event.key = 'value'"},
			},
		},
		Reporters: yamlConfig.Reporters{
			{
				Name: "test",
				Type: "telegram",
				TelegramConfig: &yamlConfig.TelegramConfig{
					Chat:   1,
					Token:  "xxx:yyy",
					Admins: []int64{123},
				},
			},
		},
		Subscriptions: yamlConfig.Subscriptions{
			{},
		},
	}
	require.Error(t, config.Validate())
}

func TestYamlConfigChainSubscriptionChainNotFound(t *testing.T) {
	t.Parallel()

	config := yamlConfig.YamlConfig{
		Chains: yamlConfig.Chains{
			{
				Name:            "chain",
				ChainID:         "chain-id",
				TendermintNodes: []string{"node"},
				APINodes:        []string{"node"},
				Queries:         []string{"event.key = 'value'"},
			},
		},
		Reporters: yamlConfig.Reporters{
			{
				Name: "test",
				Type: "telegram",
				TelegramConfig: &yamlConfig.TelegramConfig{
					Chat:   1,
					Token:  "xxx:yyy",
					Admins: []int64{123},
				},
			},
		},
		Subscriptions: yamlConfig.Subscriptions{
			{
				Name:     "name",
				Reporter: "reporter",
				ChainSubscriptions: yamlConfig.ChainSubscriptions{
					{Chain: "nonexistent"},
				},
			},
		},
	}
	require.Error(t, config.Validate())
}

func TestYamlConfigSubscriptionReporterNotFound(t *testing.T) {
	t.Parallel()

	config := yamlConfig.YamlConfig{
		Chains: yamlConfig.Chains{
			{
				Name:            "chain",
				ChainID:         "chain-id",
				TendermintNodes: []string{"node"},
				APINodes:        []string{"node"},
				Queries:         []string{"event.key = 'value'"},
			},
		},
		Reporters: yamlConfig.Reporters{
			{
				Name: "test",
				Type: "telegram",
				TelegramConfig: &yamlConfig.TelegramConfig{
					Chat:   1,
					Token:  "xxx:yyy",
					Admins: []int64{123},
				},
			},
		},
		Subscriptions: yamlConfig.Subscriptions{
			{
				Name:     "name",
				Reporter: "nonexistent",
				ChainSubscriptions: yamlConfig.ChainSubscriptions{
					{Chain: "chain"},
				},
			},
		},
	}
	require.Error(t, config.Validate())
}

func TestYamlConfigValid(t *testing.T) {
	t.Parallel()

	config := yamlConfig.YamlConfig{
		Chains: yamlConfig.Chains{
			{
				Name:            "chain",
				ChainID:         "chain-id",
				TendermintNodes: []string{"node"},
				APINodes:        []string{"node"},
				Queries:         []string{"event.key = 'value'"},
			},
		},
		Reporters: yamlConfig.Reporters{
			{
				Name: "test",
				Type: "telegram",
				TelegramConfig: &yamlConfig.TelegramConfig{
					Chat:   1,
					Token:  "xxx:yyy",
					Admins: []int64{123},
				},
			},
		},
		Subscriptions: yamlConfig.Subscriptions{
			{
				Name:     "name",
				Reporter: "test",
				ChainSubscriptions: yamlConfig.ChainSubscriptions{
					{Chain: "chain"},
				},
			},
		},
	}
	require.NoError(t, config.Validate())
}
