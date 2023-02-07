package price_fetchers

import configTypes "main/pkg/config/types"

type PriceFetcher interface {
	GetPrice(*configTypes.DenomInfo) (float64, error)
}
