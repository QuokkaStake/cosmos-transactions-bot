package types

type DenomInfo struct {
	Denom             string
	DenomExponent     int
	DisplayDenom      string
	CoingeckoCurrency string
}

func (d *DenomInfo) DisplayWarnings(chain *Chain) []DisplayWarning {
	var warnings []DisplayWarning

	if d.CoingeckoCurrency == "" {
		warnings = append(warnings, DisplayWarning{
			Keys: map[string]string{
				"chain": chain.Name,
			},
			Text: "No denoms set, prices in USD won't be displayed.",
		})
	}

	return warnings
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
