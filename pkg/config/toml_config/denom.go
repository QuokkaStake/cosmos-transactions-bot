package toml_config

import (
	"fmt"
	"main/pkg/config/types"
)

type DenomInfo struct {
	Denom             string `toml:"denom"`
	DisplayDenom      string `default:""                toml:"display-denom"`
	DenomCoefficient  int64  `default:"1000000"         toml:"denom-coefficient"`
	CoingeckoCurrency string `toml:"coingecko-currency"`
}

func (d *DenomInfo) Validate() error {
	if d.Denom == "" {
		return fmt.Errorf("denom is not set")
	}

	if d.DisplayDenom == "" {
		return fmt.Errorf("display denom is not set")
	}

	return nil
}

type DenomInfos []*DenomInfo

func (d DenomInfos) ToAppConfigDenomInfos() types.DenomInfos {
	denomInfos := make(types.DenomInfos, len(d))
	for index, info := range d {
		denomInfos[index] = &types.DenomInfo{
			Denom:             info.Denom,
			DisplayDenom:      info.DisplayDenom,
			DenomCoefficient:  info.DenomCoefficient,
			CoingeckoCurrency: info.CoingeckoCurrency,
		}
	}

	return denomInfos
}

func TomlConfigDenomsFrom(d types.DenomInfos) DenomInfos {
	denomInfos := make(DenomInfos, len(d))
	for index, info := range d {
		denomInfos[index] = &DenomInfo{
			Denom:             info.Denom,
			DisplayDenom:      info.DisplayDenom,
			DenomCoefficient:  info.DenomCoefficient,
			CoingeckoCurrency: info.CoingeckoCurrency,
		}
	}

	return denomInfos
}
