package data_fetcher

import (
	"main/assets"
	aliasManagerPkg "main/pkg/alias_manager"
	configPkg "main/pkg/config"
	"main/pkg/config/types"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"
	"main/pkg/metrics"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestDataFetcherFetchRemoteChainIdChainNotFound(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{Name: "chain"},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	data, fetched := dataFetcher.GetIbcRemoteChainID("chain-id", "channel", "port")
	require.False(t, fetched)
	require.Empty(t, data)
}

func TestDataFetcherFetchRemoteChainIdCachedOk(t *testing.T) {
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

	dataFetcher.Cache.Set("chain_channel_channel_port_port", "remote-chain")

	data, fetched := dataFetcher.GetIbcRemoteChainID("chain-id", "channel", "port")
	require.True(t, fetched)
	require.Equal(t, "remote-chain", data)
}

func TestDataFetcherFetchRemoteChainIdCachedNotOk(t *testing.T) {
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

	dataFetcher.Cache.Set("chain_channel_channel_port_port", nil)

	data, fetched := dataFetcher.GetIbcRemoteChainID("chain-id", "channel", "port")
	require.False(t, fetched)
	require.Empty(t, data)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestDataFetcherFetchRemoteChainIdIbcChannelQueryFailed(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{Name: "chain", ChainID: "chain-id", APINodes: []string{"https://example.com"}},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	data, fetched := dataFetcher.GetIbcRemoteChainID("chain-id", "channel", "port")
	require.False(t, fetched)
	require.Empty(t, data)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestDataFetcherFetchRemoteChainIdMultihop(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/ibc/core/channel/v1/channels/channel/ports/port",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("ibc-channel-multihop.json")),
	)

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{Name: "chain", ChainID: "chain-id", APINodes: []string{"https://example.com"}},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	data, fetched := dataFetcher.GetIbcRemoteChainID("chain-id", "channel", "port")
	require.False(t, fetched)
	require.Empty(t, data)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestDataFetcherFetchRemoteChainIdIbcClientStateQueryFailed(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/ibc/core/channel/v1/channels/channel/ports/port",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("ibc-channel.json")),
	)

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{Name: "chain", ChainID: "chain-id", APINodes: []string{"https://example.com"}},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	data, fetched := dataFetcher.GetIbcRemoteChainID("chain-id", "channel", "port")
	require.False(t, fetched)
	require.Empty(t, data)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestDataFetcherFetchRemoteChainIdOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/ibc/core/channel/v1/channels/channel/ports/port",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("ibc-channel.json")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/ibc/core/connection/v1/connections/connection-5/client_state",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("ibc-client-state.json")),
	)

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{Name: "chain", ChainID: "chain-id", APINodes: []string{"https://example.com"}},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	data, fetched := dataFetcher.GetIbcRemoteChainID("chain-id", "channel", "port")
	require.True(t, fetched)
	require.Equal(t, "denis-fadeev-chain", data)
}
