package data_fetcher

import "main/pkg/types/responses"

func (f *DataFetcher) GetCosmosDirectoryChains() (responses.CosmosDirectoryChains, bool) {
	keyName := "cosmos_directory_chains"

	if cachedChains, cachedChainsPresent := f.Cache.Get(keyName); cachedChainsPresent {
		if cachedChainsParsed, ok := cachedChains.(responses.CosmosDirectoryChains); ok {
			return cachedChainsParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached chains to responses.CosmosDirectoryChains")
		return nil, false
	}

	notCachedChainsList, err := f.CosmosDirectoryClient.GetAllChains()
	if err != nil {
		f.Logger.Error().Msg("Error fetching chains list")
		return nil, false
	}

	f.Cache.Set(keyName, notCachedChainsList)
	return notCachedChainsList, true
}
