package price_fetchers

import (
	configTypes "main/pkg/config/types"

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
