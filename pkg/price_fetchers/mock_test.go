package price_fetchers

import (
	"main/pkg/config/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockPriceFetcher(t *testing.T) {
	t.Parallel()

	fetcher := MockPriceFetcher{}
	require.Equal(t, "mock", fetcher.Name())

	prices, err := fetcher.GetPrices(types.DenomInfos{})
	require.NoError(t, err)
	require.Empty(t, prices)
}
