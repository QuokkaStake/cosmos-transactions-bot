package main

import (
	"github.com/rs/zerolog"
)

type DataFetcher struct {
	Logger       zerolog.Logger
	Cache        *Cache
	Chain        *Chain
	PriceFetcher PriceFetcher
}

func NewDataFetcher(logger *zerolog.Logger, chain *Chain) *DataFetcher {
	return &DataFetcher{
		Logger: logger.With().
			Str("component", "data_fetcher").
			Str("chain", chain.Name).
			Logger(),
		Cache:        NewCache(logger),
		PriceFetcher: GetPriceFetcher(logger, chain),
		Chain:        chain,
	}
}

func (f *DataFetcher) GetPrice() (float64, bool) {
	if f.PriceFetcher == nil {
		return 0, false
	}

	if cachedPrice, cachedPricePresent := f.Cache.Get(f.Chain.Name + "_price"); cachedPricePresent {
		if cachedPriceFloat, ok := cachedPrice.(float64); ok {
			return cachedPriceFloat, true
		}

		f.Logger.Error().Msg("Could not convert cached price to float64")
		return 0, false
	}

	notCachedPrice, err := f.PriceFetcher.GetPrice()
	if err != nil {
		f.Logger.Error().Msg("Error fetching price")
		return 0, false
	}

	f.Cache.Set(f.Chain.Name+"_price", notCachedPrice)
	return notCachedPrice, true
}
