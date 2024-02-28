package data_fetcher

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types/responses"
)

func (f *DataFetcher) GetProposal(chain *configTypes.Chain, id string) (*responses.Proposal, bool) {
	keyName := chain.Name + "_proposal_" + id

	if cachedEntry, cachedEntryPresent := f.Cache.Get(keyName); cachedEntryPresent {
		if cachedEntryParsed, ok := cachedEntry.(*responses.Proposal); ok {
			return cachedEntryParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached proposal to responses.Proposal")
		return nil, false
	}

	for _, node := range f.TendermintApiClients[chain.Name] {
		notCachedEntry, err := node.GetProposal(id)
		if err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching proposal")
			continue
		}

		f.Cache.Set(keyName, notCachedEntry)
		return notCachedEntry, true
	}

	f.Logger.Error().Msg("Could not connect to any nodes to get proposal")
	return nil, false
}
