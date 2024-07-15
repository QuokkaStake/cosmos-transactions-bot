package data_fetcher

import (
	aliasManagerPkg "main/pkg/alias_manager"
	configPkg "main/pkg/config"
	"main/pkg/config/types"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"
	"main/pkg/metrics"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDataFetcherFindChainById(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{Name: "chain", ChainID: "chain-id"},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	chain, found := dataFetcher.FindChainById("chain-id")
	require.True(t, found)
	require.NotNil(t, chain)

	chain2, found2 := dataFetcher.FindChainById("random")
	require.False(t, found2)
	require.Nil(t, chain2)
}

func TestDataFetcherFindSubscriptionByReporter(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Subscriptions: types.Subscriptions{
			{Name: "subscription", Reporter: "reporter"},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	subscription, found := dataFetcher.FindSubscriptionByReporter("reporter")
	require.True(t, found)
	require.NotNil(t, subscription)

	subscription2, found2 := dataFetcher.FindSubscriptionByReporter("random")
	require.False(t, found2)
	require.Nil(t, subscription2)
}

func TestDataFetcherFindChainsByReporter(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Subscriptions: types.Subscriptions{
			{Name: "subscription", Reporter: "reporter", ChainSubscriptions: types.ChainSubscriptions{
				{Chain: "chain"},
				{Chain: "chain2"},
			}},
		},
		Chains: types.Chains{
			{Name: "chain", ChainID: "chain-id"},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	chains := dataFetcher.FindChainsByReporter("reporter")
	require.NotEmpty(t, chains)

	chains2 := dataFetcher.FindChainsByReporter("random")
	require.Empty(t, chains2)
}
