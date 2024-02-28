package data_fetcher

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types/amount"
	"strings"
)

func (f *DataFetcher) PopulateMultichainDenomInfo(
	chainID string,
	baseDenom amount.Denom,
) (*configTypes.DenomInfo, bool) {
	// Getting a chain given we know its chain-id and its base denom.

	// 1. Trying to find local chain in config.
	if chain, chainFound := f.FindChainById(chainID); chainFound {
		if denom := chain.Denoms.Find(baseDenom.String()); denom != nil {
			return denom, true
		}
	}

	// 2. If it's an IBC denom - we need to fetch the remote chain's chain-id
	// and fetch it from that chain.
	if baseDenom.IsIbcToken() {
		ibcChainID, remoteDenom, fetched := f.GetRemoteChainIDAndDenomByIBCDenom(chainID, baseDenom)
		if !fetched {
			return nil, false
		}

		return f.PopulateMultichainDenomInfo(ibcChainID, remoteDenom)
	}

	// 3. Trying to fetch chain from cosmos.directory.
	// 3.1. Fetching the cosmos.directory chains list
	cosmosDirectoryChains, found := f.GetCosmosDirectoryChains()
	if !found {
		return nil, false
	}

	// 3.2. Finding the chain by chain-id from their response.
	remoteChain, found := cosmosDirectoryChains.FindByChainID(chainID)
	if !found {
		return nil, false
	}

	// 3.3. Finding the denom from their response
	remoteDenomInfo, err := remoteChain.GetDenomInfo(string(baseDenom))
	if err != nil {
		f.Logger.Error().
			Err(err).
			Str("chain", chainID).
			Str("denom", string(baseDenom)).
			Msg("Error parsing the remote denom from cosmos.directory")
		return nil, false
	}

	return remoteDenomInfo, true
}

func (f *DataFetcher) GetRemoteChainIDAndDenomByIBCDenom(
	chainID string,
	denom amount.Denom,
) (string, amount.Denom, bool) {
	// Gets IBC remote chain ID by IBC denom on a given chain.
	chain, found := f.FindChainById(chainID)
	if !found {
		return "", "", false
	}

	// 1. Fetching remote DenomTrace from chain this transaction/message belongs to.
	trace, found := f.GetDenomTrace(chain, string(denom))
	if !found {
		return "", "", false
	}

	pathParsed := strings.Split(trace.Path, "/")
	remoteChainID := chainID

	// 2. Traversing path.
	for len(pathParsed) > 0 {
		port, channel, pathParsedInternal := pathParsed[0], pathParsed[1], pathParsed[2:]
		pathParsed = pathParsedInternal

		// 3. Getting the chain-id of the denom on the chain it was minted.
		remoteChainIDFetched, found := f.GetIbcRemoteChainID(remoteChainID, channel, port)
		if !found {
			return "", "", false
		}
		remoteChainID = remoteChainIDFetched
	}

	return remoteChainID, amount.Denom(trace.BaseDenom), true
}
