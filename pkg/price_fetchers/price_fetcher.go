package price_fetchers

import configTypes "main/pkg/config/types"

type PriceFetcher interface {
	GetPrice(denomInfo *configTypes.DenomInfo) (float64, error)
}
