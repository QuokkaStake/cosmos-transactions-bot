package main

import (
	"github.com/rs/zerolog"
	gecko "github.com/superoo7/go-gecko/v3"
)

type PriceFetcher interface {
	GetPrice() (float64, error)
}

func GetPriceFetcher(logger *zerolog.Logger, chain *Chain) PriceFetcher {
	if chain.CoingeckoCurrency != "" {
		return NewCoingeckoPriceFetcher(logger, chain)
	}

	return nil
}

type CoingeckoPriceFetcher struct {
	Client *gecko.Client
	Chain  *Chain
	Logger zerolog.Logger
}

func NewCoingeckoPriceFetcher(logger *zerolog.Logger, chain *Chain) *CoingeckoPriceFetcher {
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
