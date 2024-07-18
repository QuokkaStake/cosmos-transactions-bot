package price_fetchers

import (
	"errors"
	"main/assets"
	"main/pkg/config"
	"main/pkg/config/types"
	loggerPkg "main/pkg/logger"
	"main/pkg/metrics"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // disabled due to httpmock usage
func TestCoingeckoQueryFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://api.coingecko.com/api/v3/simple/price?ids=cosmos&vs_currencies=usd",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)

	logger := loggerPkg.GetNopLogger()
	metricsManager := metrics.NewManager(logger, config.MetricsConfig{})
	coingecko := NewCoingeckoPriceFetcher(*logger, metricsManager)

	denomInfos := types.DenomInfos{{Denom: "atom", CoingeckoCurrency: "cosmos"}}
	currencies, err := coingecko.GetPrices(denomInfos)
	require.Error(t, err)
	require.ErrorContains(t, err, "custom error")
	require.Empty(t, currencies)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestCoingeckoQuerySuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://api.coingecko.com/api/v3/simple/price?ids=cosmos,akash-network,random&vs_currencies=usd",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("coingecko.json")),
	)

	logger := loggerPkg.GetNopLogger()
	metricsManager := metrics.NewManager(logger, config.MetricsConfig{})
	coingecko := NewCoingeckoPriceFetcher(*logger, metricsManager)

	denomInfos := types.DenomInfos{
		{Denom: "atom", CoingeckoCurrency: "cosmos"},
		{Denom: "akt", CoingeckoCurrency: "akash-network"},
		{Denom: "random", CoingeckoCurrency: "random"},
	}
	currencies, err := coingecko.GetPrices(denomInfos)
	require.NoError(t, err)
	require.Len(t, currencies, 2)
	require.NotNil(t, currencies[denomInfos[0]])
	require.NotNil(t, currencies[denomInfos[1]])
}
