package data_fetcher

import (
	"main/pkg/alias_manager"
	"main/pkg/cache"
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
	cosmosDirectoryPkg "main/pkg/cosmos_directory"
	"main/pkg/metrics"
	"main/pkg/tendermint/api"
	priceFetchers "main/pkg/types"

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

func (f *DataFetcher) GetAliasManager() *alias_manager.AliasManager {
	return f.AliasManager
}

func (f *DataFetcher) FindChainById(
	chainID string,
) (*configTypes.Chain, bool) {
	return f.Config.Chains.FindByChainID(chainID)
}
