package price_fetchers

import configTypes "main/pkg/config/types"

type PriceFetcher interface {
	GetPrices(denomInfos configTypes.DenomInfos) (map[*configTypes.DenomInfo]float64, error)
	Name() string
}
