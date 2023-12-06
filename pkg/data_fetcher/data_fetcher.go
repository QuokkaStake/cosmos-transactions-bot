package data_fetcher

import (
	"fmt"
	"main/pkg/metrics"
	amountPkg "main/pkg/types/amount"
	"strconv"

	"main/pkg/alias_manager"
	"main/pkg/cache"
	configTypes "main/pkg/config/types"
	priceFetchers "main/pkg/price_fetchers"
	"main/pkg/tendermint/api"
	QueryInfo "main/pkg/types/query_info"
	"main/pkg/types/responses"

	"github.com/rs/zerolog"
)

type DataFetcher struct {
	Logger               zerolog.Logger
	Cache                *cache.Cache
	Chain                *configTypes.Chain
	PriceFetchers        map[string]priceFetchers.PriceFetcher
	AliasManager         *alias_manager.AliasManager
	MetricsManager       *metrics.Manager
	TendermintApiClients []*api.TendermintApiClient
}

func NewDataFetcher(
	logger *zerolog.Logger,
	chain *configTypes.Chain,
	aliasManager *alias_manager.AliasManager,
	metricsManager *metrics.Manager,
) *DataFetcher {
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
		PriceFetchers:        map[string]priceFetchers.PriceFetcher{},
		Chain:                chain,
		TendermintApiClients: tendermintApiClients,
		AliasManager:         aliasManager,
		MetricsManager:       metricsManager,
	}
}

func (f *DataFetcher) GetPriceFetcher(info *configTypes.DenomInfo) priceFetchers.PriceFetcher {
	if info.CoingeckoCurrency != "" {
		if fetcher, ok := f.PriceFetchers[priceFetchers.CoingeckoPriceFetcherName]; ok {
			return fetcher
		}

		f.PriceFetchers[priceFetchers.CoingeckoPriceFetcherName] = priceFetchers.NewCoingeckoPriceFetcher(f.Logger)
		return f.PriceFetchers[priceFetchers.CoingeckoPriceFetcherName]
	}

	return nil
}

func (f *DataFetcher) GetDenomPriceKey(denomInfo *configTypes.DenomInfo) string {
	return fmt.Sprintf("%s_price_%s", f.Chain.Name, denomInfo.Denom)
}
func (f *DataFetcher) MaybeGetCachedPrice(denomInfo *configTypes.DenomInfo) (float64, bool) {
	cacheKey := f.GetDenomPriceKey(denomInfo)

	if cachedPrice, cachedPricePresent := f.Cache.Get(cacheKey); cachedPricePresent {
		if cachedPriceFloat, ok := cachedPrice.(float64); ok {
			return cachedPriceFloat, true
		}

		f.Logger.Error().Msg("Could not convert cached price to float64")
		return 0, false
	}

	return 0, false
}

func (f *DataFetcher) SetCachedPrice(denomInfo *configTypes.DenomInfo, notCachedPrice float64) {
	cacheKey := f.GetDenomPriceKey(denomInfo)
	f.Cache.Set(cacheKey, notCachedPrice)
}

func (f *DataFetcher) PopulateAmount(amount *amountPkg.Amount) {
	f.PopulateAmounts(amountPkg.Amounts{amount})
}

func (f *DataFetcher) PopulateAmounts(amounts amountPkg.Amounts) {
	denomsToQueryByPriceFetcher := make(map[string]configTypes.DenomInfos)

	// 1. Getting cached prices.
	for _, amount := range amounts {
		denomInfo := f.Chain.Denoms.Find(amount.BaseDenom)
		if denomInfo == nil {
			continue
		}

		amount.ConvertDenom(denomInfo.DisplayDenom, denomInfo.DenomCoefficient)

		// If we've found cached price, then using it.
		if price, cached := f.MaybeGetCachedPrice(denomInfo); cached {
			amount.AddUSDPrice(price)
			continue
		}

		// Otherwise, we try to figure out what price fetcher to use
		// and put it into a map to query it all at once.
		priceFetcher := f.GetPriceFetcher(denomInfo)
		if priceFetcher == nil {
			continue
		}

		if _, ok := denomsToQueryByPriceFetcher[priceFetcher.Name()]; !ok {
			denomsToQueryByPriceFetcher[priceFetcher.Name()] = make(configTypes.DenomInfos, 0)
		}

		denomsToQueryByPriceFetcher[priceFetcher.Name()] = append(
			denomsToQueryByPriceFetcher[priceFetcher.Name()],
			denomInfo,
		)
	}

	// 2. If we do not need to fetch any prices from price fetcher (e.g. no prices here
	// or all prices are taken from cache), then we do not need to do anything else.
	if len(denomsToQueryByPriceFetcher) == 0 {
		return
	}

	uncachedPrices := make(map[string]float64, 0)

	// 3. Fetching all prices by price fetcher.
	for priceFetcherKey, priceFetcher := range f.PriceFetchers {
		pricesToFetch, ok := denomsToQueryByPriceFetcher[priceFetcherKey]
		if !ok {
			continue
		}

		// Actually fetching prices.
		prices, err := priceFetcher.GetPrices(pricesToFetch)

		if err != nil {
			continue
		}

		// Saving it to cache
		for denomInfo, price := range prices {
			f.SetCachedPrice(denomInfo, price)

			uncachedPrices[denomInfo.Denom] = price
		}
	}

	// 4. Converting USD amounts for newly fetched prices.
	for _, amount := range amounts {
		uncachedPrice, ok := uncachedPrices[amount.BaseDenom]
		if !ok {
			continue
		}

		amount.AddUSDPrice(uncachedPrice)
	}
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
		notCachedValidator, err, queryInfo := node.GetValidator(address)
		f.MetricsManager.LogTendermintQuery(f.Chain.Name, queryInfo, QueryInfo.QueryTypeValidator)
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
		notCachedValidator, err, queryInfo := node.GetDelegatorsRewardsAtBlock(delegator, validator, block-1)
		f.MetricsManager.LogTendermintQuery(f.Chain.Name, queryInfo, QueryInfo.QueryTypeRewards)
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
		notCachedEntry, err, queryInfo := node.GetValidatorCommissionAtBlock(validator, block-1)
		f.MetricsManager.LogTendermintQuery(f.Chain.Name, queryInfo, QueryInfo.QueryTypeCommission)
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
		notCachedEntry, err, queryInfo := node.GetProposal(id)
		f.MetricsManager.LogTendermintQuery(f.Chain.Name, queryInfo, QueryInfo.QueryTypeProposal)
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
		notCachedEntry, err, queryInfo := node.GetStakingParams()
		f.MetricsManager.LogTendermintQuery(f.Chain.Name, queryInfo, QueryInfo.QueryTypeStakingParams)

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

func (f *DataFetcher) GetAliasManager() *alias_manager.AliasManager {
	return f.AliasManager
}

func (f *DataFetcher) GetChain() *configTypes.Chain {
	return f.Chain
}
