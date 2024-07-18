package data_fetcher

import (
	"main/assets"
	aliasManagerPkg "main/pkg/alias_manager"
	configPkg "main/pkg/config"
	"main/pkg/config/types"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"
	"main/pkg/metrics"
	priceFetchers "main/pkg/price_fetchers"
	amountPkg "main/pkg/types/amount"
	"math/big"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestPopulateAmountCachedNotConvertedNotFound(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
				Denoms:  types.DenomInfos{},
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	amount := &amountPkg.Amount{BaseDenom: "uatom", Value: big.NewFloat(1230000)}

	dataFetcher.PopulateAmount(config.Chains[0].ChainID, amount)
	require.Equal(t, "1230000", amount.Value.String())
	require.Equal(t, "uatom", amount.BaseDenom.String())
	require.Nil(t, amount.PriceUSD)
}

func TestPopulateAmountCachedNotConvertedWithoutPriceFetcher(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
				Denoms: types.DenomInfos{
					{Denom: "uatom", DisplayDenom: "atom", DenomExponent: 6},
				},
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	amount := &amountPkg.Amount{BaseDenom: "uatom", Denom: "uatom", Value: big.NewFloat(1230000)}

	dataFetcher.PopulateAmount(config.Chains[0].ChainID, amount)
	require.Equal(t, "1.23", amount.Value.String())
	require.Equal(t, "uatom", amount.BaseDenom.String())
	require.Equal(t, "atom", amount.Denom.String())
	require.Nil(t, amount.PriceUSD)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestPopulateAmountCachedConvertedWithCachedPrice(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
				Denoms: types.DenomInfos{
					{Denom: "uatom", DisplayDenom: "atom", DenomExponent: 6, CoingeckoCurrency: "cosmos"},
				},
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain-id_price_uatom", float64(10))

	amount := &amountPkg.Amount{BaseDenom: "uatom", Denom: "uatom", Value: big.NewFloat(1230000)}

	dataFetcher.PopulateAmount(config.Chains[0].ChainID, amount)
	require.Equal(t, "1.23", amount.Value.String())
	require.Equal(t, "uatom", amount.BaseDenom.String())
	require.Equal(t, "atom", amount.Denom.String())
	require.NotNil(t, amount.PriceUSD)
	require.Equal(t, "12.3", amount.PriceUSD.String())
}

//nolint:paralleltest // disabled due to httpmock usage
func TestPopulateAmountCachedNotConvertedWithCachedPriceInvalid(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
				Denoms: types.DenomInfos{
					{Denom: "uatom", DisplayDenom: "atom", DenomExponent: 6, CoingeckoCurrency: "cosmos"},
				},
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain-id_price_uatom", nil)

	amount := &amountPkg.Amount{BaseDenom: "uatom", Denom: "uatom", Value: big.NewFloat(1230000)}

	dataFetcher.PopulateAmount(config.Chains[0].ChainID, amount)
	require.Equal(t, "1.23", amount.Value.String())
	require.Equal(t, "uatom", amount.BaseDenom.String())
	require.Equal(t, "atom", amount.Denom.String())
	require.Nil(t, amount.PriceUSD)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestPopulateAmountsFetched(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
				Denoms: types.DenomInfos{
					{Denom: "uatom", DisplayDenom: "atom", DenomExponent: 6, CoingeckoCurrency: "cosmos"},
				},
			},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetDefaultLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.PriceFetchers[priceFetchers.MockPriceFetcherName] = &priceFetchers.MockPriceFetcher{}

	httpmock.RegisterResponder(
		"GET",
		"https://api.coingecko.com/api/v3/simple/price?ids=cosmos,cosmos&vs_currencies=usd",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("coingecko.json")),
	)

	amounts := amountPkg.Amounts{
		{BaseDenom: "uatom", Denom: "uatom", Value: big.NewFloat(1230000)},
		{BaseDenom: "uatom", Denom: "uatom", Value: big.NewFloat(4560000)},
	}

	dataFetcher.PopulateAmounts(config.Chains[0].ChainID, amounts)

	require.Equal(t, "1.23", amounts[0].Value.String())
	require.Equal(t, "uatom", amounts[0].BaseDenom.String())
	require.Equal(t, "atom", amounts[0].Denom.String())
	require.NotNil(t, amounts[0].PriceUSD)
	require.Equal(t, "8.1057", amounts[0].PriceUSD.String())

	require.Equal(t, "4.56", amounts[1].Value.String())
	require.Equal(t, "uatom", amounts[1].BaseDenom.String())
	require.Equal(t, "atom", amounts[1].Denom.String())
	require.NotNil(t, amounts[1].PriceUSD)
	require.Equal(t, "30.0504", amounts[1].PriceUSD.String())
}
