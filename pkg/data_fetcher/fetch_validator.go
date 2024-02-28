package data_fetcher

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types/responses"
)

func (f *DataFetcher) GetValidator(chain *configTypes.Chain, address string) (*responses.Validator, bool) {
	keyName := chain.Name + "_validator_" + address

	if cachedValidator, cachedValidatorPresent := f.Cache.Get(keyName); cachedValidatorPresent {
		if cachedValidatorParsed, ok := cachedValidator.(*responses.Validator); ok {
			return cachedValidatorParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached validator to *stakingTypes.Validator")
		return nil, false
	}

	for _, node := range f.TendermintApiClients[chain.Name] {
		notCachedValidator, err := node.GetValidator(address)
		if err != nil {
			f.Logger.Error().Msg("Error fetching validator")
			continue
		}

		f.Cache.Set(keyName, notCachedValidator)
		return notCachedValidator, true
	}

	f.Logger.Error().Msg("Could not connect to any nodes to get a validator")
	return nil, false
}
