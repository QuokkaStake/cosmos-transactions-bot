package data_fetcher

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types/responses"
	"strconv"
)

func (f *DataFetcher) GetCommissionAtBlock(
	chain *configTypes.Chain,
	validator string,
	block int64,
) ([]responses.Commission, bool) {
	keyName := chain.Name + "_commission_" + validator + "_" + strconv.FormatInt(block, 10)

	if cachedCommission, cachedCommissionPresent := f.Cache.Get(keyName); cachedCommissionPresent {
		if cachedCommissionParsed, ok := cachedCommission.([]responses.Commission); ok {
			return cachedCommissionParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached commission to []responses.Commission")
		return []responses.Commission{}, false
	}

	for _, node := range f.TendermintApiClients[chain.Name] {
		notCachedEntry, err := node.GetValidatorCommissionAtBlock(validator, block-1)
		if err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching commission")
			continue
		}

		f.Cache.Set(keyName, notCachedEntry)
		return notCachedEntry, true
	}

	f.Logger.Error().Msg("Could not connect to any nodes to get commission")
	return []responses.Commission{}, false
}
