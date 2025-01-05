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

func TestDataFetcherPopulateWalletExplorerNotPresent(t *testing.T) {
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

	dataFetcher.PopulateWallet(config.Chains[0], wallet, "subscription")
	require.Empty(t, wallet.Href)
	require.Empty(t, wallet.Title)
}

func TestDataFetcherPopulateWalletPresent(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:     "chain",
				Explorer: &types.Explorer{WalletLinkPattern: "link %s"},
			},
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
	dataFetcher.PopulateWallet(config.Chains[0], wallet, "subscription")

	require.Equal(t, "link address", wallet.Href)
	require.Equal(t, "alias", wallet.Title)
}

func TestDataFetcherPopulateMultichainWalletNoChannelOrPort(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:     "chain",
				Explorer: &types.Explorer{WalletLinkPattern: "link %s"},
			},
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
	dataFetcher.PopulateMultichainWallet(config.Chains[0], "", "", wallet, "subscription")

	require.Equal(t, "link address", wallet.Href)
	require.Equal(t, "alias", wallet.Title)
}
func TestDataFetcherPopulateMultichainWalletNoRemoteChainFetched(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:     "chain",
				Explorer: &types.Explorer{WalletLinkPattern: "link %s"},
			},
		},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_channel_channel_port_port", nil)

	err := aliasManager.Set("subscription", "chain", "address", "alias")
	require.NoError(t, err)

	wallet := &types.Link{Value: "address"}
	dataFetcher.PopulateMultichainWallet(config.Chains[0], "channel", "port", wallet, "subscription")

	require.Empty(t, wallet.Href)
	require.Empty(t, wallet.Title)
}

func TestDataFetcherPopulateMultichainWalletNoLocalChainFetched(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:     "chain",
				Explorer: &types.Explorer{WalletLinkPattern: "link %s"},
			},
		},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_channel_channel_port_port", "remote-chain")

	err := aliasManager.Set("subscription", "chain", "address", "alias")
	require.NoError(t, err)

	wallet := &types.Link{Value: "address"}
	dataFetcher.PopulateMultichainWallet(config.Chains[0], "channel", "port", wallet, "subscription")

	require.Empty(t, wallet.Href)
	require.Empty(t, wallet.Title)
}

func TestDataFetcherPopulateMultichainWalletNoLocalChainExplorer(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:     "chain",
				Explorer: &types.Explorer{WalletLinkPattern: "link %s"},
			},
			{
				Name:    "chain2",
				ChainID: "remote-chain",
			},
		},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_channel_channel_port_port", "remote-chain")

	err := aliasManager.Set("subscription", "chain", "address", "alias")
	require.NoError(t, err)

	wallet := &types.Link{Value: "address"}
	dataFetcher.PopulateMultichainWallet(config.Chains[0], "channel", "port", wallet, "subscription")

	require.Empty(t, wallet.Href)
	require.Empty(t, wallet.Title)
}

func TestDataFetcherPopulateMultichainWalletRemoteOk(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:     "chain",
				Explorer: &types.Explorer{WalletLinkPattern: "link %s"},
			},
			{
				Name:     "chain2",
				ChainID:  "remote-chain",
				Explorer: &types.Explorer{WalletLinkPattern: "another link %s"},
			},
		},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_channel_channel_port_port", "remote-chain")

	err := aliasManager.Set("subscription", "chain2", "address", "alias")
	require.NoError(t, err)

	wallet := &types.Link{Value: "address"}
	dataFetcher.PopulateMultichainWallet(config.Chains[0], "channel", "port", wallet, "subscription")

	require.Equal(t, "another link address", wallet.Href)
	require.Equal(t, "alias", wallet.Title)
}
