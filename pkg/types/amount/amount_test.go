package amount_test

import (
	"fmt"
	amountPkg "main/pkg/types/amount"
	"testing"

	sdkmath "cosmossdk.io/math"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestAmountFrom(t *testing.T) {
	t.Parallel()

	coin := cosmosTypes.Coin{Denom: "stake", Amount: sdkmath.NewInt(100)}
	amount := amountPkg.AmountFrom(coin)

	require.Equal(t, "stake", amount.Denom)
	require.Equal(t, "stake", amount.BaseDenom)
	require.Equal(t, "100.000000", fmt.Sprintf("%.6f", amount.Value))
}

func TestAmountFromStringInvalid(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	_ = amountPkg.AmountFromString("invalid", "stake")
}

func TestAmountFromStringValid(t *testing.T) {
	t.Parallel()

	amount := amountPkg.AmountFromString("100", "stake")

	require.Equal(t, "stake", amount.Denom)
	require.Equal(t, "stake", amount.BaseDenom)
	require.Equal(t, "100.000000", fmt.Sprintf("%.6f", amount.Value))
}

func TestAmountConvertDenom(t *testing.T) {
	t.Parallel()

	amount := amountPkg.AmountFromString("100000000", "ustake")
	amount.ConvertDenom("stake", 1000000)

	require.Equal(t, "stake", amount.Denom)
	require.Equal(t, "ustake", amount.BaseDenom)
	require.Equal(t, "100.000000", fmt.Sprintf("%.6f", amount.Value))
}

func TestAmountAddUsdPrice(t *testing.T) {
	t.Parallel()

	amount := amountPkg.AmountFromString("1", "stake")
	amount.AddUSDPrice(1.23)

	require.Equal(t, "stake", amount.Denom)
	require.Equal(t, "1.000000", fmt.Sprintf("%.6f", amount.Value))
	require.Equal(t, "1.230000", fmt.Sprintf("%.6f", amount.PriceUSD))
}

func TestAmountToString(t *testing.T) {
	t.Parallel()

	amount := amountPkg.AmountFromString("123.456", "stake")
	require.Equal(t, "123stake", amount.String())
}

func TestAmountIsIbcDenom(t *testing.T) {
	t.Parallel()

	amount := amountPkg.AmountFromString("123.456", "ibc/xxxxx")
	require.True(t, amount.IsIbcToken())

	amount2 := amountPkg.AmountFromString("123.456", "ustake")
	require.False(t, amount2.IsIbcToken())
}

func TestAmountsToString(t *testing.T) {
	t.Parallel()

	amount1 := amountPkg.AmountFromString("123.456", "stake")
	amount2 := amountPkg.AmountFromString("345.678", "yield")

	amounts := amountPkg.Amounts{amount1, amount2}

	require.Equal(t, "123stake,345yield", amounts.String())
}
