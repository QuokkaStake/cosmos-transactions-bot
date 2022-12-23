package price_fetchers

import (
	"github.com/rs/zerolog"
	"main/pkg/types/chains"
)

type PriceFetcher interface {
	GetPrice() (float64, error)
}

func GetPriceFetcher(logger *zerolog.Logger, chain *chains.Chain) PriceFetcher {
	if chain.CoingeckoCurrency != "" {
		return NewCoingeckoPriceFetcher(logger, chain)
	}

	return nil
}
