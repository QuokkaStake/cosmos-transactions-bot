package data_fetcher

import (
	aliasManagerPkg "main/pkg/alias_manager"
	configPkg "main/pkg/config"
	"main/pkg/config/types"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"
	"main/pkg/metrics"
	"main/pkg/types/responses"
	"testing"

	transferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"
)

func TestDataFetcherGetMultichainDenomInfoPresentLocally(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
				Denoms:  types.DenomInfos{{Denom: "udenom", DisplayDenom: "denom"}},
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	denomInfo, found := dataFetcher.PopulateMultichainDenomInfo("chain-id", "udenom")
	require.True(t, found)
	require.NotNil(t, denomInfo)
	require.Equal(t, "denom", denomInfo.DisplayDenom)
}

func TestDataFetcherGetMultichainDenomInfoIbcDenomNotFetched(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:   "chain",
				Denoms: types.DenomInfos{{Denom: "udenom", DisplayDenom: "denom"}},
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	denomInfo, found := dataFetcher.PopulateMultichainDenomInfo("chain-id", "ibc/denom")
	require.False(t, found)
	require.Nil(t, denomInfo)
}

func TestDataFetcherGetMultichainDenomInfoIbcDenomTraceFailed(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
				Denoms:  types.DenomInfos{{Denom: "udenom", DisplayDenom: "denom"}},
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_denom_trace_denom", nil)

	denomInfo, found := dataFetcher.PopulateMultichainDenomInfo("chain-id", "ibc/denom")
	require.False(t, found)
	require.Nil(t, denomInfo)
}

func TestDataFetcherGetMultichainDenomInfoIbcDenomRemoteFetchIdFailed(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
				Denoms:  types.DenomInfos{{Denom: "udenom", DisplayDenom: "denom"}},
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_denom_trace_denom", &transferTypes.DenomTrace{
		Path: "port/channel",
	})

	dataFetcher.Cache.Set("chain_channel_channel_port_port", nil)

	denomInfo, found := dataFetcher.PopulateMultichainDenomInfo("chain-id", "ibc/denom")
	require.False(t, found)
	require.Nil(t, denomInfo)
}

func TestDataFetcherGetMultichainDenomInfoIbcDenomMultihopRemoteChainNotFound(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
				Denoms:  types.DenomInfos{{Denom: "udenom", DisplayDenom: "denom"}},
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_denom_trace_denom", &transferTypes.DenomTrace{
		Path: "port/channel/port2/channel2",
	})

	dataFetcher.Cache.Set("chain_channel_channel_port_port", "remotechain")
	dataFetcher.Cache.Set("chain_channel_channel2_port_port2", "remotechain2")

	denomInfo, found := dataFetcher.PopulateMultichainDenomInfo("chain-id", "ibc/denom")
	require.False(t, found)
	require.Nil(t, denomInfo)
}

func TestDataFetcherGetMultichainDenomInfoIbcDenomRemoteChainNotFound(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
				Denoms:  types.DenomInfos{{Denom: "udenom", DisplayDenom: "denom"}},
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_denom_trace_denom", &transferTypes.DenomTrace{
		Path: "port/channel",
	})

	dataFetcher.Cache.Set("chain_channel_channel_port_port", "remotechain")

	denomInfo, found := dataFetcher.PopulateMultichainDenomInfo("chain-id", "ibc/denom")
	require.False(t, found)
	require.Nil(t, denomInfo)
}

func TestDataFetcherGetMultichainDenomInfoCosmosDirectoryFailed(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name: "chain",
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("cosmos_directory_chains", nil)

	denomInfo, found := dataFetcher.PopulateMultichainDenomInfo("chain-id", "udenom")
	require.False(t, found)
	require.Nil(t, denomInfo)
}

func TestDataFetcherGetMultichainDenomInfoCosmosDirectoryChainNotFound(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name: "chain",
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("cosmos_directory_chains", responses.CosmosDirectoryChains{})

	denomInfo, found := dataFetcher.PopulateMultichainDenomInfo("chain-id", "udenom")
	require.False(t, found)
	require.Nil(t, denomInfo)
}

func TestDataFetcherGetMultichainDenomInfoCosmosDirectoryAssetNotFound(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name: "chain",
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("cosmos_directory_chains", responses.CosmosDirectoryChains{
		{
			ChainID: "chain-id",
			Assets:  []responses.CosmosDirectoryAsset{},
		},
	})

	denomInfo, found := dataFetcher.PopulateMultichainDenomInfo("chain-id", "udenom")
	require.False(t, found)
	require.Nil(t, denomInfo)
}

func TestDataFetcherGetMultichainDenomInfoCosmosDirectoryOk(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name: "chain",
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("cosmos_directory_chains", responses.CosmosDirectoryChains{
		{
			ChainID: "chain-id",
			Assets: []responses.CosmosDirectoryAsset{
				{
					Denom:   "udenom",
					Base:    responses.CosmosDirectoryAssetDenomInfo{Denom: "udenom"},
					Display: responses.CosmosDirectoryAssetDenomInfo{Denom: "denom"},
				},
			},
		},
	})

	denomInfo, found := dataFetcher.PopulateMultichainDenomInfo("chain-id", "udenom")
	require.True(t, found)
	require.NotNil(t, denomInfo)
	require.Equal(t, "denom", denomInfo.DisplayDenom)
}
