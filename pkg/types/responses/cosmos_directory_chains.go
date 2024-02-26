package responses

import (
	"fmt"
	"main/pkg/config/types"
	"math"
)

type CosmosDirectoryChainsResponse struct {
	Chains CosmosDirectoryChains `json:"chains"`
}

type CosmosDirectoryChains []CosmosDirectoryChain

func (chains CosmosDirectoryChains) FindByChainID(chainID string) (CosmosDirectoryChain, bool) {
	for _, chain := range chains {
		if chain.ChainID == chainID {
			return chain, true
		}
	}

	return CosmosDirectoryChain{}, false
}

type CosmosDirectoryChain struct {
	Name    string                 `json:"name"`
	ChainID string                 `json:"chain_id"`
	Assets  []CosmosDirectoryAsset `json:"assets"`
}

func (chain CosmosDirectoryChain) GetDenomInfo(baseDenom string) (*types.DenomInfo, error) {
	for _, asset := range chain.Assets {
		if asset.Denom != baseDenom {
			continue
		}

		if asset.Base.Denom == "" || asset.Display.Denom == "" {
			return nil, fmt.Errorf(
				"got malformed cosmos.directory response: base.denom '%s', display.denom '%s'",
				asset.Base.Denom,
				asset.Display.Denom,
			)
		}

		return &types.DenomInfo{
			Denom:             asset.Base.Denom,
			DisplayDenom:      asset.Display.Denom,
			CoingeckoCurrency: asset.CoingeckoID,
			DenomCoefficient:  int64(math.Pow10(asset.Display.Exponent - asset.Base.Exponent)),
		}, nil
	}

	return nil, fmt.Errorf("asset is not found on chain %s\n", chain.ChainID)
}

type CosmosDirectoryAsset struct {
	Denom       string                        `json:"denom"`
	CoingeckoID string                        `json:"coingecko_id"`
	Base        CosmosDirectoryAssetDenomInfo `json:"base"`
	Display     CosmosDirectoryAssetDenomInfo `json:"display"`
}

type CosmosDirectoryAssetDenomInfo struct {
	Denom    string `json:"denom"`
	Exponent int    `json:"exponent"`
}
