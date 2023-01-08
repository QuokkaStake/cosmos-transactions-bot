package data_fetcher

import (
	"strconv"

	"main/pkg/alias_manager"
	"main/pkg/cache"
	"main/pkg/config/types"
	priceFetchers "main/pkg/price_fetchers"
	"main/pkg/tendermint/api"
	"main/pkg/types/responses"

	"github.com/rs/zerolog"
)

type DataFetcher struct {
	Logger               zerolog.Logger
	Cache                *cache.Cache
	Chain                *types.Chain
	PriceFetcher         priceFetchers.PriceFetcher
	AliasManager         *alias_manager.AliasManager
	TendermintApiClients []*api.TendermintApiClient
}

func NewDataFetcher(logger *zerolog.Logger, chain *types.Chain, aliasManager *alias_manager.AliasManager) *DataFetcher {
	tendermintApiClients := make([]*api.TendermintApiClient, len(chain.APINodes))
	for index, node := range chain.APINodes {
		tendermintApiClients[index] = api.NewTendermintApiClient(logger, node, chain)
	}

	return &DataFetcher{
		Logger: logger.With().
			Str("component", "data_fetcher").
			Str("chain", chain.Name).
			Logger(),
		Cache:                cache.NewCache(logger),
		PriceFetcher:         priceFetchers.GetPriceFetcher(logger, chain),
		Chain:                chain,
		TendermintApiClients: tendermintApiClients,
		AliasManager:         aliasManager,
	}
}

func (f *DataFetcher) GetPrice() (float64, bool) {
	if f.PriceFetcher == nil {
		return 0, false
	}

	if cachedPrice, cachedPricePresent := f.Cache.Get(f.Chain.Name + "_price"); cachedPricePresent {
		if cachedPriceFloat, ok := cachedPrice.(float64); ok {
			return cachedPriceFloat, true
		}

		f.Logger.Error().Msg("Could not convert cached price to float64")
		return 0, false
	}

	notCachedPrice, err := f.PriceFetcher.GetPrice()
	if err != nil {
		f.Logger.Error().Msg("Error fetching price")
		return 0, false
	}

	f.Cache.Set(f.Chain.Name+"_price", notCachedPrice)
	return notCachedPrice, true
}

func (f *DataFetcher) GetValidator(address string) (*responses.Validator, bool) {
	keyName := f.Chain.Name + "_validator_" + address

	if cachedValidator, cachedValidatorPresent := f.Cache.Get(keyName); cachedValidatorPresent {
		if cachedValidatorParsed, ok := cachedValidator.(*responses.Validator); ok {
			return cachedValidatorParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached validator to *stakingTypes.Validator")
		return nil, false
	}

	for _, node := range f.TendermintApiClients {
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

func (f *DataFetcher) GetRewardsAtBlock(
	delegator string,
	validator string,
	block int64,
) ([]responses.Reward, bool) {
	keyName := f.Chain.Name + "_rewards_" + delegator + "_" + validator + "_" + strconv.FormatInt(block, 10)

	if cachedRewards, cachedRewardsPresent := f.Cache.Get(keyName); cachedRewardsPresent {
		if cachedRewardsParsed, ok := cachedRewards.([]responses.Reward); ok {
			return cachedRewardsParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached rewards to []responses.Reward")
		return []responses.Reward{}, false
	}

	for _, node := range f.TendermintApiClients {
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

func (f *DataFetcher) GetCommissionAtBlock(
	validator string,
	block int64,
) ([]responses.Commission, bool) {
	keyName := f.Chain.Name + "_commission_" + validator + "_" + strconv.FormatInt(block, 10)

	if cachedCommission, cachedCommissionPresent := f.Cache.Get(keyName); cachedCommissionPresent {
		if cachedCommissionParsed, ok := cachedCommission.([]responses.Commission); ok {
			return cachedCommissionParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached commission to []responses.Commission")
		return []responses.Commission{}, false
	}

	for _, node := range f.TendermintApiClients {
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

func (f *DataFetcher) GetProposal(id string) (*responses.Proposal, bool) {
	keyName := f.Chain.Name + "_proposal_" + id

	if cachedEntry, cachedEntryPresent := f.Cache.Get(keyName); cachedEntryPresent {
		if cachedEntryParsed, ok := cachedEntry.(*responses.Proposal); ok {
			return cachedEntryParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached proposal to responses.Proposal")
		return nil, false
	}

	for _, node := range f.TendermintApiClients {
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

func (f *DataFetcher) GetStakingParams() (*responses.StakingParams, bool) {
	keyName := f.Chain.Name + "_staking_params"

	if cachedEntry, cachedEntryPresent := f.Cache.Get(keyName); cachedEntryPresent {
		if cachedEntryParsed, ok := cachedEntry.(*responses.StakingParams); ok {
			return cachedEntryParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached staking params to *responses.StakingParams")
		return nil, false
	}

	for _, node := range f.TendermintApiClients {
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
