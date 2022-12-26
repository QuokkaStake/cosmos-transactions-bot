package price_fetchers

import (
	"github.com/rs/zerolog"
	gecko "github.com/superoo7/go-gecko/v3"
	"main/pkg/config/types"
)

type CoingeckoPriceFetcher struct {
	Client *gecko.Client
	Chain  *types.Chain
	Logger zerolog.Logger
}

func NewCoingeckoPriceFetcher(logger *zerolog.Logger, chain *types.Chain) *CoingeckoPriceFetcher {
	return &CoingeckoPriceFetcher{
		Client: gecko.NewClient(nil),
		Logger: logger.With().Str("component", "coingecko_price_fetcher").Logger(),
		Chain:  chain,
	}
}

func (c *CoingeckoPriceFetcher) GetPrice() (float64, error) {
	result, err := c.Client.SimpleSinglePrice(c.Chain.CoingeckoCurrency, "usd")
	if err != nil {
		c.Logger.Error().Err(err).Msg("Could not get rate")
		return 0, err
	}

	return float64(result.MarketPrice), nil
}
