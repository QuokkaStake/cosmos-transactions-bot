package price_fetchers

import (
	"main/pkg/config/types"

	"github.com/rs/zerolog"
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
