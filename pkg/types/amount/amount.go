package amount

import (
	"fmt"
	"main/pkg/logger"
	"main/pkg/utils"
	"math/big"
	"strings"

	transferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
)

type Denom string

func (d Denom) IsIbcToken() bool {
	denomSplit := strings.Split(string(d), "/")
	return len(denomSplit) == 2 && denomSplit[0] == transferTypes.DenomPrefix
}

func (d Denom) String() string {
	return string(d)
}

type Amount struct {
	Value     *big.Float
	Denom     Denom
	BaseDenom Denom
	PriceUSD  *big.Float
}

func AmountFrom(coin cosmosTypes.Coin) *Amount {
	return &Amount{
		Value:     new(big.Float).SetInt(coin.Amount.BigInt()),
		Denom:     Denom(coin.Denom),
		BaseDenom: Denom(coin.Denom),
	}
}

func AmountFromString(amount string, denom string) *Amount {
	parsedAmount, ok := new(big.Float).SetString(amount)
	if !ok {
		logger.GetDefaultLogger().Panic().Str("value", amount).Msg("Could not parse string as big.Float")
	}

	return &Amount{
		Value:     parsedAmount,
		Denom:     Denom(denom),
		BaseDenom: Denom(denom),
	}
}

func (a *Amount) ConvertDenom(displayDenom string, denomCoefficient int64) {
	denomCoefficientBigFloat := new(big.Float).SetInt64(denomCoefficient)
	a.Value = new(big.Float).Quo(a.Value, denomCoefficientBigFloat)
	a.BaseDenom = a.Denom
	a.Denom = Denom(displayDenom)
}

func (a *Amount) AddUSDPrice(usdPrice float64) {
	tokenPriceBigFloat := new(big.Float).Set(a.Value)
	amountValueBigFloat := new(big.Float).SetFloat64(usdPrice)
	a.PriceUSD = new(big.Float).Mul(tokenPriceBigFloat, amountValueBigFloat)
}

func (a Amount) String() string {
	value, _ := a.Value.Int(nil)
	return fmt.Sprintf("%d%s", value, a.Denom)
}

type Amounts []*Amount

func (a Amounts) String() string {
	return strings.Join(utils.Map(a, func(a *Amount) string {
		return a.String()
	}), ",")
}
