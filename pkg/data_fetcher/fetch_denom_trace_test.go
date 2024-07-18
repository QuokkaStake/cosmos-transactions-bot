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

	transferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestDataFetcherFetchDenomTraceInvalidDenom(t *testing.T) {
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

	data, fetched := dataFetcher.GetDenomTrace(config.Chains[0], "invalid")
	require.False(t, fetched)
	require.Nil(t, data)
}

func TestDataFetcherFetchDenomTraceCachedOk(t *testing.T) {
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

	dataFetcher.Cache.Set("chain_denom_trace_denom", &transferTypes.DenomTrace{
		Path: "path",
	})

	data, fetched := dataFetcher.GetDenomTrace(config.Chains[0], "ibc/denom")
	require.True(t, fetched)
	require.NotNil(t, data)
	require.Equal(t, "path", data.Path)
}

func TestDataFetcherFetchDenomTraceCachedNotOk(t *testing.T) {
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

	dataFetcher.Cache.Set("chain_denom_trace_denom", nil)

	data, fetched := dataFetcher.GetDenomTrace(config.Chains[0], "ibc/denom")
	require.False(t, fetched)
	require.Nil(t, data)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestDataFetcherFetchDenomTraceAllQueriesFailed(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{Name: "chain", APINodes: []string{"https://example.com"}},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	data, fetched := dataFetcher.GetDenomTrace(config.Chains[0], "ibc/denom")
	require.False(t, fetched)
	require.Nil(t, data)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestDataFetcherFetchDenomTraceSuccessfullyFetched(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/ibc/apps/transfer/v1/denom_traces/denom",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("denom-trace.json")),
	)

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{Name: "chain", APINodes: []string{"https://example.com"}},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	data, fetched := dataFetcher.GetDenomTrace(config.Chains[0], "ibc/denom")
	require.True(t, fetched)
	require.NotNil(t, data)
	require.Equal(t, "untrn", data.BaseDenom)
}
