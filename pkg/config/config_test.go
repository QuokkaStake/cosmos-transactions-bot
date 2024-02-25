package config_test

import (
	"main/assets"
	configPkg "main/pkg/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfigErrorReading(t *testing.T) {
	t.Parallel()

	config, err := configPkg.GetConfig("config.toml", assets.AssetsFS)

	require.Error(t, err)
	require.Nil(t, config)
}

func TestLoadConfigInvalidConfig(t *testing.T) {
	t.Parallel()

	config, err := configPkg.GetConfig("invalid-timezone.toml", assets.AssetsFS)

	require.Error(t, err)
	require.Nil(t, config)
}

func TestLoadConfigInvalidToml(t *testing.T) {
	t.Parallel()

	config, err := configPkg.GetConfig("invalid-toml.toml", assets.AssetsFS)

	require.Error(t, err)
	require.Nil(t, config)
}

func TestLoadConfigValid(t *testing.T) {
	t.Parallel()

	config, err := configPkg.GetConfig("valid.toml", assets.AssetsFS)

	require.NoError(t, err)
	require.NotNil(t, config)
}

func TestConfigDisplayWarnings(t *testing.T) {
	t.Parallel()

	config, err := configPkg.GetConfig("valid.toml", assets.AssetsFS)

	require.NoError(t, err)
	require.NotNil(t, config)

	warnings := config.DisplayWarnings()
	require.Empty(t, warnings)
}

func TestConfigDisplayAsToml(t *testing.T) {
	t.Parallel()

	config, err := configPkg.GetConfig("valid.toml", assets.AssetsFS)
	require.NoError(t, err)

	tomlConfig := config.ToTomlConfig()
	configAgain := configPkg.FromTomlConfig(tomlConfig)

	require.EqualValues(t, config.LogConfig, configAgain.LogConfig)
	require.EqualValues(t, config.AliasesPath, configAgain.AliasesPath)
	require.EqualValues(t, config.Metrics, configAgain.Metrics)
	require.EqualValues(t, config.Timezone, configAgain.Timezone)

	require.Equal(t, len(config.Chains), len(configAgain.Chains))
	for index := range config.Chains {
		configChain := config.Chains[index]
		configChainAgain := configAgain.Chains[index]

		require.EqualValues(t, configChain.Name, configChainAgain.Name)
		require.EqualValues(t, configChain.PrettyName, configChainAgain.PrettyName)
		require.EqualValues(t, configChain.ChainID, configChainAgain.ChainID)
		require.EqualValues(t, configChain.TendermintNodes, configChainAgain.TendermintNodes)
		require.EqualValues(t, configChain.APINodes, configChainAgain.APINodes)
		require.EqualValues(t, configChain.Explorer, configChainAgain.Explorer)
		require.EqualValues(t, configChain.SupportedExplorer, configChainAgain.SupportedExplorer)
		require.EqualValues(t, configChain.Denoms, configChainAgain.Denoms)

		require.Equal(t, len(configChain.Queries), len(configChainAgain.Queries))
		for queryIndex := range configChain.Queries {
			require.EqualValues(
				t,
				configChain.Queries[queryIndex].String(),
				configChainAgain.Queries[queryIndex].String(),
			)
		}
	}

	require.Equal(t, len(config.Subscriptions), len(configAgain.Subscriptions))
	for index := range config.Subscriptions {
		configSubscription := config.Subscriptions[index]
		configSubscriptionAgain := configAgain.Subscriptions[index]

		require.EqualValues(t, configSubscription.Name, configSubscriptionAgain.Name)
		require.EqualValues(t, configSubscription.Reporter, configSubscriptionAgain.Reporter)

		require.Equal(
			t,
			len(configSubscription.ChainSubscriptions),
			len(configSubscriptionAgain.ChainSubscriptions),
		)

		for chainSubscriptionIndex := range configSubscription.ChainSubscriptions {
			configChainSubscription := configSubscription.ChainSubscriptions[chainSubscriptionIndex]
			configChainSubscriptionAgain := configSubscriptionAgain.ChainSubscriptions[chainSubscriptionIndex]

			require.EqualValues(t, configChainSubscription.Chain, configChainSubscriptionAgain.Chain)
			require.EqualValues(t, configChainSubscription.LogNodeErrors, configChainSubscriptionAgain.LogNodeErrors)
			require.EqualValues(t, configChainSubscription.LogUnparsedMessages, configChainSubscriptionAgain.LogUnparsedMessages)
			require.EqualValues(t, configChainSubscription.LogUnknownMessages, configChainSubscriptionAgain.LogUnknownMessages)
			require.EqualValues(t, configChainSubscription.LogFailedTransactions, configChainSubscriptionAgain.LogFailedTransactions)
			require.EqualValues(t, configChainSubscription.FilterInternalMessages, configChainSubscriptionAgain.FilterInternalMessages)

			require.Equal(
				t,
				len(configChainSubscription.Filters),
				len(configChainSubscriptionAgain.Filters),
			)

			for filterIndex := range configChainSubscription.Filters {
				require.Equal(
					t,
					len(configChainSubscription.Filters[filterIndex].String()),
					len(configChainSubscriptionAgain.Filters[filterIndex].String()),
				)
			}
		}
	}

	require.Equal(t, len(config.Reporters), len(configAgain.Reporters))
	for index := range config.Reporters {
		require.EqualValues(t, config.Reporters[index], configAgain.Reporters[index])
	}
}

func TestGetConfigAsString(t *testing.T) {
	t.Parallel()

	config, err := configPkg.GetConfig("valid.toml", assets.AssetsFS)

	require.NoError(t, err)
	require.NotNil(t, config)

	configString := config.GetConfigAsString()
	require.NotEmpty(t, configString)
}
