package data_fetcher

import (
	"fmt"
	configPkg "main/pkg/config"
	"main/pkg/metrics"
	amountPkg "main/pkg/types/amount"
	"strconv"
	"strings"

	transferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

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
	Config               *configPkg.AppConfig
	PriceFetchers        map[string]priceFetchers.PriceFetcher
	AliasManager         *alias_manager.AliasManager
	MetricsManager       *metrics.Manager
	TendermintApiClients map[string][]*api.TendermintApiClient
}

func NewDataFetcher(
	logger *zerolog.Logger,
	config *configPkg.AppConfig,
	aliasManager *alias_manager.AliasManager,
	metricsManager *metrics.Manager,
) *DataFetcher {
	tendermintApiClients := make(map[string][]*api.TendermintApiClient, len(config.Chains))
	for _, chain := range config.Chains {
		tendermintApiClients[chain.Name] = make([]*api.TendermintApiClient, len(chain.APINodes))
		for index, node := range chain.APINodes {
			tendermintApiClients[chain.Name][index] = api.NewTendermintApiClient(logger, node, chain)
		}
	}

	return &DataFetcher{
		Logger: logger.With().
			Str("component", "data_fetcher").
			Logger(),
		Cache:                cache.NewCache(logger),
		PriceFetchers:        map[string]priceFetchers.PriceFetcher{},
		Config:               config,
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

func (f *DataFetcher) GetDenomPriceKey(
	chain *configTypes.Chain,
	denomInfo *configTypes.DenomInfo,
) string {
	return fmt.Sprintf("%s_price_%s", chain.Name, denomInfo.Denom)
}
func (f *DataFetcher) MaybeGetCachedPrice(
	chain *configTypes.Chain,
	denomInfo *configTypes.DenomInfo,
) (float64, bool) {
	cacheKey := f.GetDenomPriceKey(chain, denomInfo)

	if cachedPrice, cachedPricePresent := f.Cache.Get(cacheKey); cachedPricePresent {
		if cachedPriceFloat, ok := cachedPrice.(float64); ok {
			return cachedPriceFloat, true
		}

		f.Logger.Error().Msg("Could not convert cached price to float64")
		return 0, false
	}

	return 0, false
}

func (f *DataFetcher) SetCachedPrice(
	chain *configTypes.Chain,
	denomInfo *configTypes.DenomInfo,
	notCachedPrice float64,
) {
	cacheKey := f.GetDenomPriceKey(chain, denomInfo)
	f.Cache.Set(cacheKey, notCachedPrice)
}

func (f *DataFetcher) PopulateAmount(chain *configTypes.Chain, amount *amountPkg.Amount) {
	f.PopulateAmounts(chain, amountPkg.Amounts{amount})
}

func (f *DataFetcher) PopulateAmounts(chain *configTypes.Chain, amounts amountPkg.Amounts) {
	denomsToQueryByPriceFetcher := make(map[string]configTypes.DenomInfos)

	// 1. Getting cached prices.
	for _, amount := range amounts {
		denomInfo := chain.Denoms.Find(amount.BaseDenom)
		if denomInfo == nil {
			continue
		}

		amount.ConvertDenom(denomInfo.DisplayDenom, denomInfo.DenomCoefficient)

		// If we've found cached price, then using it.
		if price, cached := f.MaybeGetCachedPrice(chain, denomInfo); cached {
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

	uncachedPrices := make(map[string]float64)

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
			f.SetCachedPrice(chain, denomInfo, price)

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
		notCachedValidator, err, queryInfo := node.GetValidator(address)
		f.MetricsManager.LogTendermintQuery(chain.Name, queryInfo, QueryInfo.QueryTypeValidator)
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
		notCachedValidator, err, queryInfo := node.GetDelegatorsRewardsAtBlock(delegator, validator, block-1)
		f.MetricsManager.LogTendermintQuery(chain.Name, queryInfo, QueryInfo.QueryTypeRewards)
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
		notCachedEntry, err, queryInfo := node.GetValidatorCommissionAtBlock(validator, block-1)
		f.MetricsManager.LogTendermintQuery(chain.Name, queryInfo, QueryInfo.QueryTypeCommission)
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
		notCachedEntry, err, queryInfo := node.GetProposal(id)
		f.MetricsManager.LogTendermintQuery(chain.Name, queryInfo, QueryInfo.QueryTypeProposal)
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
		notCachedEntry, err, queryInfo := node.GetStakingParams()
		f.MetricsManager.LogTendermintQuery(chain.Name, queryInfo, QueryInfo.QueryTypeStakingParams)

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

func (f *DataFetcher) GetIbcRemoteChainID(
	chain *configTypes.Chain,
	channel string,
	port string,
) (string, bool) {
	keyName := chain.Name + "_channel_" + channel + "_port_" + port

	if cachedEntry, cachedEntryPresent := f.Cache.Get(keyName); cachedEntryPresent {
		if cachedEntryParsed, ok := cachedEntry.(string); ok {
			return cachedEntryParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached IBC channel to string")
		return "", false
	}

	var (
		ibcChannel     *responses.IbcChannel
		ibcClientState *responses.IbcIdentifiedClientState
	)

	for _, node := range f.TendermintApiClients[chain.Name] {
		ibcChannelResponse, err, queryInfo := node.GetIbcChannel(channel, port)
		f.MetricsManager.LogTendermintQuery(chain.Name, queryInfo, QueryInfo.QueryTypeIbcChannel)
		if err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching IBC channel")
			continue
		}

		ibcChannel = ibcChannelResponse
		break
	}

	if ibcChannel == nil {
		f.Logger.Error().Msg("Could not connect to any nodes to get IBC channel")
		return "", false
	}

	if len(ibcChannel.ConnectionHops) != 1 {
		f.Logger.Error().
			Int("len", len(ibcChannel.ConnectionHops)).
			Msg("Could not connect to any nodes to get IBC channel")
		return "", false
	}

	for _, node := range f.TendermintApiClients[chain.Name] {
		ibcChannelClientStateResponse, err, queryInfo := node.GetIbcConnectionClientState(ibcChannel.ConnectionHops[0])
		f.MetricsManager.LogTendermintQuery(chain.Name, queryInfo, QueryInfo.QueryTypeIbcConnectionClientState)
		if err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching IBC client state")
			continue
		}

		ibcClientState = ibcChannelClientStateResponse
		break
	}

	if ibcClientState == nil {
		f.Logger.Error().Msg("Could not connect to any nodes to get IBC client state")
		return "", false
	}

	f.Cache.Set(keyName, ibcClientState.ClientState.ChainId)
	return ibcClientState.ClientState.ChainId, true
}

func (f *DataFetcher) GetDenomTrace(
	chain *configTypes.Chain,
	denom string,
) (*transferTypes.DenomTrace, bool) {
	denomSplit := strings.Split(denom, "/")
	if len(denomSplit) != 2 || denomSplit[0] != transferTypes.DenomPrefix {
		f.Logger.Error().Msg("Invalid IBC prefix provided")
		return nil, false
	}

	denomHash := denomSplit[1]

	keyName := chain.Name + "_denom_trace_" + denom

	if cachedEntry, cachedEntryPresent := f.Cache.Get(keyName); cachedEntryPresent {
		if cachedEntryParsed, ok := cachedEntry.(*transferTypes.DenomTrace); ok {
			return cachedEntryParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached staking params to *types.DenomTrace")
		return nil, false
	}

	for _, node := range f.TendermintApiClients[chain.Name] {
		notCachedEntry, err, queryInfo := node.GetIbcDenomTrace(denomHash)
		f.MetricsManager.LogTendermintQuery(chain.Name, queryInfo, QueryInfo.QueryTypeIbcDenomTrace)

		if err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching IBC denom trace")
			continue
		}

		f.Cache.Set(keyName, notCachedEntry)
		return notCachedEntry, true
	}

	f.Logger.Error().Msg("Could not connect to any nodes to get IBC denom trace")
	return nil, false
}

func (f *DataFetcher) GetAliasManager() *alias_manager.AliasManager {
	return f.AliasManager
}

func (f *DataFetcher) FindChainById(
	chainID string,
) (*configTypes.Chain, bool) {
	chain := f.Config.Chains.FindByChainID(chainID)
	if chain == nil {
		return nil, false
	}

	return chain, true
}
