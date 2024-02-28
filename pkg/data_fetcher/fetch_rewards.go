package data_fetcher

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types/responses"
	"strconv"
)

func (f *DataFetcher) GetRewardsAtBlock(
	chain *configTypes.Chain,
	delegator string,
	validator string,
	block int64,
) ([]responses.Reward, bool) {
	keyName := chain.Name + "_rewards_" + delegator + "_" + validator + "_" + strconv.FormatInt(block, 10)

	if cachedRewards, cachedRewardsPresent := f.Cache.Get(keyName); cachedRewardsPresent {
		if cachedRewardsParsed, ok := cachedRewards.([]responses.Reward); ok {
			return cachedRewardsParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached rewards to []responses.Reward")
		return []responses.Reward{}, false
	}

	for _, node := range f.TendermintApiClients[chain.Name] {
		notCachedValidator, err := node.GetDelegatorsRewardsAtBlock(delegator, validator, block-1)
		if err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching rewards")
			continue
		}

		f.Cache.Set(keyName, notCachedValidator)
		return notCachedValidator, true
	}

	f.Logger.Error().Msg("Could not connect to any nodes to get rewards list")
	return []responses.Reward{}, false
}
