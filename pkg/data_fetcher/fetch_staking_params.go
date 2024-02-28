package data_fetcher

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types/responses"
)

func (f *DataFetcher) GetStakingParams(chain *configTypes.Chain) (*responses.StakingParams, bool) {
	keyName := chain.Name + "_staking_params"

	if cachedEntry, cachedEntryPresent := f.Cache.Get(keyName); cachedEntryPresent {
		if cachedEntryParsed, ok := cachedEntry.(*responses.StakingParams); ok {
			return cachedEntryParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached staking params to *responses.StakingParams")
		return nil, false
	}

	for _, node := range f.TendermintApiClients[chain.Name] {
		notCachedEntry, err := node.GetStakingParams()

		if err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching staking params")
			continue
		}

		f.Cache.Set(keyName, notCachedEntry)
		return notCachedEntry, true
	}

	f.Logger.Error().Msg("Could not connect to any nodes to get staking params")
	return nil, false
}
