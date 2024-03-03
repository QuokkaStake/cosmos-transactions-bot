package responses_test

import (
	"main/pkg/types/responses"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCosmosDirectoryChainsFind(t *testing.T) {
	t.Parallel()

	chains := responses.CosmosDirectoryChains{
		{ChainID: "chain"},
	}

	chain1, found1 := chains.FindByChainID("chain")
	require.NotNil(t, chain1)
	require.True(t, found1)

	_, found2 := chains.FindByChainID("chain2")
	require.False(t, found2)
}

func TestMalformedResponse(t *testing.T) {
	t.Parallel()

	chain := responses.CosmosDirectoryChain{
		ChainID: "chain",
		Assets: []responses.CosmosDirectoryAsset{
			{
				Denom:   "denom",
				Base:    responses.CosmosDirectoryAssetDenomInfo{},
				Display: responses.CosmosDirectoryAssetDenomInfo{},
			},
		},
	}

	denom, err := chain.GetDenomInfo("denom")
	require.Error(t, err)
	require.Nil(t, denom)
}

func TestDenomNotFound(t *testing.T) {
	t.Parallel()

	chain := responses.CosmosDirectoryChain{
		ChainID: "chain",
		Assets: []responses.CosmosDirectoryAsset{
			{
				Base:    responses.CosmosDirectoryAssetDenomInfo{},
				Display: responses.CosmosDirectoryAssetDenomInfo{},
			},
		},
	}

	denom, err := chain.GetDenomInfo("denom")
	require.Error(t, err)
	require.Nil(t, denom)
}

func TestDenomFound(t *testing.T) {
	t.Parallel()

	chain := responses.CosmosDirectoryChain{
		ChainID: "chain",
		Assets: []responses.CosmosDirectoryAsset{
			{
				Denom:       "udenom",
				Base:        responses.CosmosDirectoryAssetDenomInfo{Denom: "udenom", Exponent: 3},
				Display:     responses.CosmosDirectoryAssetDenomInfo{Denom: "denom", Exponent: 9},
				CoingeckoID: "coingecko",
			},
		},
	}

	denom, err := chain.GetDenomInfo("udenom")
	require.NoError(t, err)
	require.NotNil(t, denom)
	require.Equal(t, "udenom", denom.Denom)
	require.Equal(t, "denom", denom.DisplayDenom)
	require.Equal(t, "coingecko", denom.CoingeckoCurrency)
	require.Equal(t, 6, denom.DenomExponent)
}
