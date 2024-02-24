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

type Amount struct {
	Value     *big.Float
	Denom     string
	BaseDenom string
	PriceUSD  *big.Float
}

func AmountFrom(coin cosmosTypes.Coin) *Amount {
	return &Amount{
		Value:     new(big.Float).SetInt(coin.Amount.BigInt()),
		Denom:     coin.Denom,
		BaseDenom: coin.Denom,
	}
}

func AmountFromString(amount string, denom string) *Amount {
	parsedAmount, ok := new(big.Float).SetString(amount)
	if !ok {
		logger.GetDefaultLogger().Fatal().Str("value", amount).Msg("Could not parse string as big.Float")
	}

	return &Amount{
		Value:     parsedAmount,
		Denom:     denom,
		BaseDenom: denom,
	}
}

func (a *Amount) ConvertDenom(displayDenom string, denomCoefficient int64) {
	denomCoefficientBigFloat := new(big.Float).SetInt64(denomCoefficient)
	a.Value = new(big.Float).Quo(a.Value, denomCoefficientBigFloat)
	a.BaseDenom = a.Denom
	a.Denom = displayDenom
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

func (a Amount) IsIbcToken() bool {
	denomSplit := strings.Split(a.Denom, "/")
	return len(denomSplit) == 2 && denomSplit[0] == transferTypes.DenomPrefix
}

type Amounts []*Amount

func (a Amounts) String() string {
	return strings.Join(utils.Map(a, func(a *Amount) string {
		return a.String()
	}), ",")
}
