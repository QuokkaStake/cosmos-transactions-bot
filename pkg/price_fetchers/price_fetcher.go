package price_fetchers

import (
	"github.com/rs/zerolog"
	"main/pkg/config/types"
)

type PriceFetcher interface {
	GetPrice() (float64, error)
}

func GetPriceFetcher(logger *zerolog.Logger, chain *types.Chain) PriceFetcher {
	if chain.CoingeckoCurrency != "" {
		return NewCoingeckoPriceFetcher(logger, chain)
	}

	return nil
}
