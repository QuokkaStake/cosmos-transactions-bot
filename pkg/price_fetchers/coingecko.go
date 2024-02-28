package price_fetchers

import (
	"fmt"
	configTypes "main/pkg/config/types"
	"main/pkg/http"
	"main/pkg/metrics"
	"main/pkg/types/query_info"
	"main/pkg/utils"
	"strings"

	"github.com/rs/zerolog"
)

type CoingeckoPriceFetcher struct {
	Client         *http.Client
	MetricsManager *metrics.Manager
	Logger         zerolog.Logger
}

func NewCoingeckoPriceFetcher(logger zerolog.Logger, metricsManager *metrics.Manager) *CoingeckoPriceFetcher {
	return &CoingeckoPriceFetcher{
		Client:         http.NewClient(&logger, "https://api.coingecko.com", "coingecko"),
		MetricsManager: metricsManager,
		Logger:         logger.With().Str("component", "coingecko_price_fetcher").Logger(),
	}
}

func (c *CoingeckoPriceFetcher) GetPrices(denomInfos configTypes.DenomInfos) (map[*configTypes.DenomInfo]float64, error) {
	currenciesToFetch := utils.Map(denomInfos, func(denomInfo *configTypes.DenomInfo) string {
		return denomInfo.CoingeckoCurrency
	})

	var coingeckoResponse map[string]map[string]float64
	err, queryInfo := c.Client.Get(
		fmt.Sprintf(
			"/api/v3/simple/price?ids=%s&vs_currencies=%s",
			strings.Join(currenciesToFetch, ","),
			CoingeckoBaseCurrency,
		),
		&coingeckoResponse,
	)
	c.MetricsManager.LogQuery("coingecko", queryInfo, query_info.QueryTypePrices)
	if err != nil {
		c.Logger.Error().
			Err(err).
			Strs("currencies", currenciesToFetch).
			Msg("Could not get rates, probably rate-limiting")
	}

	result := make(map[*configTypes.DenomInfo]float64)

	for _, denomInfo := range denomInfos {
		coinPrice, ok := coingeckoResponse[denomInfo.CoingeckoCurrency]
		if !ok {
			continue
		}

		usdCoinPrice, ok := coinPrice[CoingeckoBaseCurrency]
		if !ok {
			result[denomInfo] = 0
		}

		result[denomInfo] = usdCoinPrice
	}

	return result, nil
}

func (c *CoingeckoPriceFetcher) Name() string {
	return CoingeckoPriceFetcherName
}
