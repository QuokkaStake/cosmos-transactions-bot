package data_fetcher

import (
	"main/assets"
	aliasManagerPkg "main/pkg/alias_manager"
	configPkg "main/pkg/config"
	"main/pkg/config/types"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"
	"main/pkg/metrics"
	"main/pkg/types/responses"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestDataFetcherFetchCommissionCachedOk(t *testing.T) {
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

	dataFetcher.Cache.Set("chain_commission_validator_100", []responses.Commission{
		{Amount: "100", Denom: "ustake"},
	})

	data, fetched := dataFetcher.GetCommissionAtBlock(config.Chains[0], "validator", 100)
	require.True(t, fetched)
	require.Len(t, data, 1)
	require.Equal(t, "100", data[0].Amount)
	require.Equal(t, "ustake", data[0].Denom)
}

func TestDataFetcherFetchCommissionCachedNotOk(t *testing.T) {
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

	dataFetcher.Cache.Set("chain_commission_validator_100", nil)

	data, fetched := dataFetcher.GetCommissionAtBlock(config.Chains[0], "validator", 100)
	require.False(t, fetched)
	require.Empty(t, data)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestDataFetcherFetchCommissionAllQueriesFailed(t *testing.T) {
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

	data, fetched := dataFetcher.GetCommissionAtBlock(config.Chains[0], "validator", 100)
	require.False(t, fetched)
	require.Empty(t, data)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestDataFetcherFetchCommissionSuccessfullyFetched(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/distribution/v1beta1/validators/validator/commission",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("commission.json")),
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

	data, fetched := dataFetcher.GetCommissionAtBlock(config.Chains[0], "validator", 100)
	require.True(t, fetched)
	require.Len(t, data, 1)
	require.Equal(t, "12345", data[0].Amount)
	require.Equal(t, "uatom", data[0].Denom)
}
