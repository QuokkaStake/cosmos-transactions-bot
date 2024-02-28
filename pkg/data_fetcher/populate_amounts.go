package data_fetcher

import (
	"fmt"
	configTypes "main/pkg/config/types"
	priceFetchers "main/pkg/price_fetchers"
	"main/pkg/types"
	amountPkg "main/pkg/types/amount"
)

func (f *DataFetcher) GetPriceFetcher(info *configTypes.DenomInfo) types.PriceFetcher {
	if info.CoingeckoCurrency != "" {
		if fetcher, ok := f.PriceFetchers[priceFetchers.CoingeckoPriceFetcherName]; ok {
			return fetcher
		}

		f.PriceFetchers[priceFetchers.CoingeckoPriceFetcherName] = priceFetchers.NewCoingeckoPriceFetcher(
			f.Logger,
			f.MetricsManager,
		)
		return f.PriceFetchers[priceFetchers.CoingeckoPriceFetcherName]
	}

	return nil
}

func (f *DataFetcher) GetDenomPriceKey(
	chainID string,
	denomInfo *configTypes.DenomInfo,
) string {
	return fmt.Sprintf("%s_price_%s", chainID, denomInfo.Denom)
}
func (f *DataFetcher) MaybeGetCachedPrice(
	chainID string,
	denomInfo *configTypes.DenomInfo,
) (float64, bool) {
	cacheKey := f.GetDenomPriceKey(chainID, denomInfo)

	if cachedPrice, cachedPricePresent := f.Cache.Get(cacheKey); cachedPricePresent {
		if cachedPriceFloat, ok := cachedPrice.(float64); ok {
			return cachedPriceFloat, true
		}

		f.Logger.Error().Msg("Could not convert cached price to float64")
		return 0, false
	}

	return 0, false
}

func (f *DataFetcher) SetCachedPrice(
	chainID string,
	denomInfo *configTypes.DenomInfo,
	notCachedPrice float64,
) {
	cacheKey := f.GetDenomPriceKey(chainID, denomInfo)
	f.Cache.Set(cacheKey, notCachedPrice)
}

func (f *DataFetcher) PopulateAmount(chainID string, amount *amountPkg.Amount) {
	f.PopulateAmounts(chainID, amountPkg.Amounts{amount})
}

func (f *DataFetcher) PopulateAmounts(chainID string, amounts amountPkg.Amounts) {
	denomsToQueryByPriceFetcher := make(map[string]configTypes.DenomInfos)

	// 1. Getting cached prices.
	for _, amount := range amounts {
		denomInfo, found := f.PopulateMultichainDenomInfo(chainID, amount.BaseDenom)
		if !found {
			f.Logger.Warn().
				Str("chain", chainID).
				Str("denom", amount.Denom.String()).
				Msg("Could not fetch denom info")
			continue
		}

		f.Logger.Debug().
			Str("chain", chainID).
			Str("denom", amount.Denom.String()).
			Str("display_denom", denomInfo.DisplayDenom).
			Int64("coefficient", denomInfo.DenomCoefficient).
			Msg("Fetched denom for chain")

		amount.ConvertDenom(denomInfo.DisplayDenom, denomInfo.DenomCoefficient)

		// If we've found cached price, then using it.
		if price, cached := f.MaybeGetCachedPrice(chainID, denomInfo); cached {
			if price != 0 {
				amount.AddUSDPrice(price)
			}
			continue
		}

		// Otherwise, we try to figure out what price fetcher to use
		// and put it into a map to query it all at once.
		priceFetcher := f.GetPriceFetcher(denomInfo)
		if priceFetcher == nil {
			continue
		}

		if _, ok := denomsToQueryByPriceFetcher[priceFetcher.Name()]; !ok {
			denomsToQueryByPriceFetcher[priceFetcher.Name()] = make(configTypes.DenomInfos, 0)
		}

		denomsToQueryByPriceFetcher[priceFetcher.Name()] = append(
			denomsToQueryByPriceFetcher[priceFetcher.Name()],
			denomInfo,
		)
	}

	// 2. If we do not need to fetch any prices from price fetcher (e.g. no prices here
	// or all prices are taken from cache), then we do not need to do anything else.
	if len(denomsToQueryByPriceFetcher) == 0 {
		return
	}

	uncachedPrices := make(map[string]float64)

	// 3. Fetching all prices by price fetcher.
	for priceFetcherKey, priceFetcher := range f.PriceFetchers {
		pricesToFetch, ok := denomsToQueryByPriceFetcher[priceFetcherKey]
		if !ok {
			continue
		}

		// Actually fetching prices.
		prices, err := priceFetcher.GetPrices(pricesToFetch)

		if err != nil {
			continue
		}

		// Saving it to cache
		for denomInfo, price := range prices {
			f.SetCachedPrice(chainID, denomInfo, price)

			uncachedPrices[denomInfo.Denom] = price
		}
	}

	// 4. Converting USD amounts for newly fetched prices.
	for _, amount := range amounts {
		uncachedPrice, ok := uncachedPrices[amount.BaseDenom.String()]
		if !ok {
			continue
		}

		if uncachedPrice != 0 {
			amount.AddUSDPrice(uncachedPrice)
		}
	}
}
