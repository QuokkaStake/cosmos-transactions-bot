package yaml_config

import (
	"fmt"
	"main/pkg/config/types"
)

type DenomInfo struct {
	Denom             string `yaml:"denom"`
	DisplayDenom      string `default:""                yaml:"display-denom"`
	DenomExponent     int    `default:"6"               yaml:"denom-exponent"`
	CoingeckoCurrency string `yaml:"coingecko-currency"`
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
			DenomExponent:     info.DenomExponent,
			CoingeckoCurrency: info.CoingeckoCurrency,
		}
	}

	return denomInfos
}

func YamlConfigDenomsFrom(d types.DenomInfos) DenomInfos {
	denomInfos := make(DenomInfos, len(d))
	for index, info := range d {
		denomInfos[index] = &DenomInfo{
			Denom:             info.Denom,
			DisplayDenom:      info.DisplayDenom,
			DenomExponent:     info.DenomExponent,
			CoingeckoCurrency: info.CoingeckoCurrency,
		}
	}

	return denomInfos
}
