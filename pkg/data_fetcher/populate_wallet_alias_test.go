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

func TestDataFetcherPopulateWalletAliasNotPresent(t *testing.T) {
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

	wallet := &types.Link{Value: "address"}

	dataFetcher.PopulateWalletAlias(config.Chains[0], wallet, "subscription")
	require.Empty(t, wallet.Title)
}

func TestDataFetcherPopulateWalletAliasPresent(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{Name: "chain"},
		},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	err := aliasManager.Set("subscription", "chain", "address", "alias")
	require.NoError(t, err)

	wallet := &types.Link{Value: "address"}

	dataFetcher.PopulateWalletAlias(config.Chains[0], wallet, "subscription")
	require.Equal(t, "alias", wallet.Title)
}
