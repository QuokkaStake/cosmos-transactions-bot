package amount

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAmountFrom(t *testing.T) {
	t.Parallel()

	coin := cosmosTypes.Coin{Denom: "stake", Amount: sdkmath.NewInt(100)}
	amount := AmountFrom(coin)

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

	_ = AmountFromString("invalid", "stake")
}

func TestAmountFromStringValid(t *testing.T) {
	t.Parallel()

	amount := AmountFromString("100", "stake")

	require.Equal(t, "stake", amount.Denom)
	require.Equal(t, "stake", amount.BaseDenom)
	require.Equal(t, "100.000000", fmt.Sprintf("%.6f", amount.Value))
}

func TestAmountConvertDenom(t *testing.T) {
	t.Parallel()

	amount := AmountFromString("100000000", "ustake")
	amount.ConvertDenom("stake", 1000000)

	require.Equal(t, "stake", amount.Denom)
	require.Equal(t, "ustake", amount.BaseDenom)
	require.Equal(t, "100.000000", fmt.Sprintf("%.6f", amount.Value))
}

func TestAmountAddUsdPrice(t *testing.T) {
	t.Parallel()

	amount := AmountFromString("1", "stake")
	amount.AddUSDPrice(1.23)

	require.Equal(t, "stake", amount.Denom)
	require.Equal(t, "1.000000", fmt.Sprintf("%.6f", amount.Value))
	require.Equal(t, "1.230000", fmt.Sprintf("%.6f", amount.PriceUSD))
}

func TestAmountToString(t *testing.T) {
	t.Parallel()

	amount := AmountFromString("123.456", "stake")
	require.Equal(t, "123stake", amount.String())
}

func TestAmountIsIbcDenom(t *testing.T) {
	t.Parallel()

	amount := AmountFromString("123.456", "ibc/xxxxx")
	require.True(t, amount.IsIbcToken())

	amount2 := AmountFromString("123.456", "ustake")
	require.False(t, amount2.IsIbcToken())
}

func TestAmountsToString(t *testing.T) {
	t.Parallel()

	amount1 := AmountFromString("123.456", "stake")
	amount2 := AmountFromString("345.678", "yield")

	amounts := Amounts{amount1, amount2}

	require.Equal(t, "123stake,345yield", amounts.String())
}
