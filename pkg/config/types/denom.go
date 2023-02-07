package types

import "github.com/rs/zerolog"

type DenomInfo struct {
	Denom             string
	DenomCoefficient  int64
	DisplayDenom      string
	CoingeckoCurrency string
}

func (d *DenomInfo) DisplayWarnings(chain *Chain, logger *zerolog.Logger) {
	if d.CoingeckoCurrency == "" {
		logger.Warn().
			Str("chain", chain.Name).
			Str("denom", d.Denom).
			Msg("Coingecko currency not set, denoms won't be displayed correctly.")
	}
}

type DenomInfos []*DenomInfo

func (d DenomInfos) Find(denom string) *DenomInfo {
	for _, info := range d {
		if denom == info.Denom {
			return info
		}
	}

	return nil
}
