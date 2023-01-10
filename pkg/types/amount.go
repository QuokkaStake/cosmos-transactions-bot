package types

import (
	"fmt"
	"main/pkg/utils"
	"math"
	"strings"
)

type Amount struct {
	Value    float64
	Denom    string
	PriceUSD float64
}

func (a Amount) String() string {
	return fmt.Sprintf("%d%s", int64(math.Round(a.Value)), a.Denom)
}

type Amounts []*Amount

func (a Amounts) String() string {
	return strings.Join(utils.Map(a, func(a *Amount) string {
		return a.String()
	}), ",")
}
