package price_fetchers

import (
	configTypes "main/pkg/config/types"
	"main/pkg/utils"

	"github.com/rs/zerolog"
	gecko "github.com/superoo7/go-gecko/v3"
)

type CoingeckoPriceFetcher struct {
	Client *gecko.Client
	Logger zerolog.Logger
}

func NewCoingeckoPriceFetcher(logger zerolog.Logger) *CoingeckoPriceFetcher {
	return &CoingeckoPriceFetcher{
		Client: gecko.NewClient(nil),
		Logger: logger.With().Str("component", "coingecko_price_fetcher").Logger(),
	}
}

func (c *CoingeckoPriceFetcher) GetPrice(denomInfo *configTypes.DenomInfo) (float64, error) {
	result, err := c.Client.SimpleSinglePrice(denomInfo.CoingeckoCurrency, "usd")
	if err != nil {
		c.Logger.Error().Err(err).Msg("Could not get rate")
		return 0, err
	}

	return float64(result.MarketPrice), nil
}

func (c *CoingeckoPriceFetcher) GetPrices(denomInfos configTypes.DenomInfos) (map[string]float64, error) {
	currenciesToFetch := utils.Map(denomInfos, func(denomInfo *configTypes.DenomInfo) string {
		return denomInfo.CoingeckoCurrency
	})

	pricesRaw, err := c.Client.SimplePrice(
		currenciesToFetch,
		[]string{"usd"},
	)

	if err != nil || pricesRaw == nil {
		c.Logger.Error().Err(err).Msg("Could not get rates")
		return make(map[string]float64, 0), err
	}

	prices := *pricesRaw

	result := make(map[string]float64, 0)

	for _, denomInfo := range denomInfos {
		coinPrice, ok := prices[denomInfo.CoingeckoCurrency]
		if !ok {
			continue
		}

		usdCoinPrice, ok := coinPrice["usd"]
		if !ok {
			continue
		}

		result[denomInfo.Denom] = float64(usdCoinPrice)
	}

	return result, nil
}

func (c *CoingeckoPriceFetcher) Name() string {
	return CoingeckoPriceFetcherName
}
