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
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestDataFetcherFetchStakingParamsCachedOk(t *testing.T) {
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

	dataFetcher.Cache.Set("chain_staking_params", &responses.StakingParams{
		UnbondingTime: responses.Duration{Duration: 15 * time.Second},
	})

	data, fetched := dataFetcher.GetStakingParams(config.Chains[0])
	require.True(t, fetched)
	require.NotNil(t, data)
	require.Equal(t, "15s", data.UnbondingTime.Duration.String())
}

func TestDataFetcherFetchStakingParamsCachedNotOk(t *testing.T) {
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

	dataFetcher.Cache.Set("chain_staking_params", nil)

	data, fetched := dataFetcher.GetStakingParams(config.Chains[0])
	require.False(t, fetched)
	require.Nil(t, data)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestDataFetcherFetchStakingParamsAllQueriesFailed(t *testing.T) {
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

	data, fetched := dataFetcher.GetStakingParams(config.Chains[0])
	require.False(t, fetched)
	require.Nil(t, data)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestDataFetcherFetchStakingParamsSuccessfullyFetched(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/params",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("staking-params.json")),
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

	data, fetched := dataFetcher.GetStakingParams(config.Chains[0])
	require.True(t, fetched)
	require.NotNil(t, data)
	require.Equal(t, "504h0m0s", data.UnbondingTime.Duration.String())
}
