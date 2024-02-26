package cosmos_directory

import (
	"main/pkg/http"
	"main/pkg/metrics"
	"main/pkg/types/query_info"
	"main/pkg/types/responses"

	"github.com/rs/zerolog"
)

type Client struct {
	Logger         zerolog.Logger
	Client         *http.Client
	MetricsManager *metrics.Manager
}

func NewClient(logger *zerolog.Logger, metricsManager *metrics.Manager) *Client {
	return &Client{
		Logger: logger.With().
			Str("component", "cosmos_directory_client").
			Logger(),
		Client:         http.NewClient(logger, "https://chains.cosmos.directory", "cosmos.directory"),
		MetricsManager: metricsManager,
	}
}

func (c *Client) GetAllChains() (responses.CosmosDirectoryChains, error) {
	var response *responses.CosmosDirectoryChainsResponse
	err, queryInfo := c.Client.Get("/", &response)
	c.MetricsManager.LogQuery("cosmos.directory", queryInfo, query_info.QueryTypeChainsList)

	if err != nil {
		return nil, err
	}

	return response.Chains, nil
}
