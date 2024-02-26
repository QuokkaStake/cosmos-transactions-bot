package data_fetcher

import (
	"fmt"
	configPkg "main/pkg/config"
	cosmosDirectoryPkg "main/pkg/cosmos_directory"
	"main/pkg/metrics"
	amountPkg "main/pkg/types/amount"
	"strings"

	"main/pkg/alias_manager"
	"main/pkg/cache"
	configTypes "main/pkg/config/types"
	priceFetchers "main/pkg/price_fetchers"
	"main/pkg/tendermint/api"
	"main/pkg/types/responses"

	"github.com/rs/zerolog"
)

type DataFetcher struct {
	Logger                zerolog.Logger
	Cache                 *cache.Cache
	Config                *configPkg.AppConfig
	PriceFetchers         map[string]priceFetchers.PriceFetcher
	AliasManager          *alias_manager.AliasManager
	MetricsManager        *metrics.Manager
	CosmosDirectoryClient *cosmosDirectoryPkg.Client

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
			tendermintApiClients[chain.Name][index] = api.NewTendermintApiClient(
				logger,
				node,
				chain,
				metricsManager,
			)
		}
	}

	return &DataFetcher{
		Logger: logger.With().
			Str("component", "data_fetcher").
			Logger(),
		Cache:                 cache.NewCache(logger),
		PriceFetchers:         map[string]priceFetchers.PriceFetcher{},
		Config:                config,
		TendermintApiClients:  tendermintApiClients,
		AliasManager:          aliasManager,
		MetricsManager:        metricsManager,
		CosmosDirectoryClient: cosmosDirectoryPkg.NewClient(logger, metricsManager),
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
		denomInfo, found := f.GetChainDenom(chain, amount)
		if !found {
			f.Logger.Warn().
				Str("chain", chain.Name).
				Str("denom", amount.Denom.String()).
				Msg("Could not fetch denom info")
			continue
		}

		f.Logger.Debug().
			Str("chain", chain.Name).
			Str("denom", amount.Denom.String()).
			Str("display_denom", denomInfo.DisplayDenom).
			Int64("coefficient", denomInfo.DenomCoefficient).
			Msg("Fetched denom for chain")

		amount.ConvertDenom(denomInfo.DisplayDenom, denomInfo.DenomCoefficient)

		// If we've found cached price, then using it.
		if price, cached := f.MaybeGetCachedPrice(chain, denomInfo); cached {
			if price != 0 {
				amount.AddUSDPrice(price)
			}
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
		uncachedPrice, ok := uncachedPrices[amount.BaseDenom.String()]
		if !ok {
			continue
		}

		if uncachedPrice != 0 {
			amount.AddUSDPrice(uncachedPrice)
		}
	}
}

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

func (f *DataFetcher) GetCosmosDirectoryDenom(chainID string, denom string) (*configTypes.DenomInfo, bool) {
	// 1. Fetching the cosmos.directory chains list
	cosmosDirectoryChains, found := f.GetCosmosDirectoryChains()
	if !found {
		return nil, false
	}

	// 2. Finding the chain by chain-id from their response.
	remoteChain, found := cosmosDirectoryChains.FindByChainID(chainID)
	if !found {
		return nil, false
	}

	// 3. Finding the denom from their response
	remoteDenomInfo, err := remoteChain.GetDenomInfo(denom)
	if err != nil {
		f.Logger.Error().
			Err(err).
			Str("chain", chainID).
			Str("denom", denom).
			Msg("Error parsing the remote denom from cosmos.directory")
		return nil, false
	}

	return remoteDenomInfo, true
}

func (f *DataFetcher) GetChainDenom(
	chain *configTypes.Chain,
	amount *amountPkg.Amount,
) (*configTypes.DenomInfo, bool) {
	// Getting the denom for that chain and denom.

	// 1. Trying to find it in local chain config.
	denomInfo := chain.Denoms.Find(amount.BaseDenom.String())
	if denomInfo != nil {
		return denomInfo, true
	}

	// 2. If it's the IBC token: trying to fetch its info from remote chain
	// (as in, query the remote chain, take chain-id from it, and then
	// take denom from cosmos.directory or local config).

	if amount.Denom.IsIbcToken() {
		return f.MaybeFetchMultichainDenom(chain, amount.BaseDenom.String())
	}

	// 3. Okay, it's not an IBC denom, we couldn't find it in the local config,
	// then we fetch the denom info from cosmos.directory.
	return f.GetCosmosDirectoryDenom(chain.ChainID, amount.BaseDenom.String())
}

func (f *DataFetcher) MaybeFetchMultichainDenom(
	chain *configTypes.Chain,
	denom string,
) (*configTypes.DenomInfo, bool) {
	// Fetching multichain denoms (as in, denoms from chains declared locally).

	// 1. Fetching remote DenomTrace from chain this transaction/message belongs to.
	trace, found := f.GetDenomTrace(chain, denom)
	if !found {
		return nil, false
	}

	// 2. Split port and channel. Multi-hop transfers are not supported.
	pathParsed := strings.Split(trace.Path, "/")
	if len(pathParsed) != 2 {
		f.Logger.Warn().
			Str("chain", chain.Name).
			Str("path", trace.Path).
			Msg("Multi-hop transfers are not yet supported.")
		return nil, false
	}

	// 3. Getting the chain-id of the denom on the chain it was minted.
	originalChainId, found := f.GetIbcRemoteChainID(chain, pathParsed[1], pathParsed[0])
	if !found {
		return nil, false
	}

	// 4. Trying to find this chain by chain-id in our local config.
	if remoteChain, remoteChainFound := f.FindChainById(originalChainId); remoteChainFound {
		if remoteDenom := remoteChain.Denoms.Find(trace.BaseDenom); remoteDenom != nil {
			return remoteDenom, true
		}
	}

	// 5. Everything failed, trying to fetch the denom from cosmos.directory
	return f.GetCosmosDirectoryDenom(originalChainId, trace.BaseDenom)
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
