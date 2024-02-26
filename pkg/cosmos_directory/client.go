package cosmos_directory

import (
	"main/pkg/http"
	"main/pkg/types/query_info"
	"main/pkg/types/responses"

	"github.com/rs/zerolog"
)

type Client struct {
	Logger zerolog.Logger
	Client *http.Client
}

func NewClient(logger *zerolog.Logger) *Client {
	return &Client{
		Logger: logger.With().
			Str("component", "cosmos_directory_client").
			Logger(),
		Client: http.NewClient(logger, "https://chains.cosmos.directory", "cosmos.directory"),
	}
}

func (c *Client) GetAllChains() (responses.CosmosDirectoryChains, error, query_info.QueryInfo) {
	var response *responses.CosmosDirectoryChainsResponse
	err, queryInfo := c.Client.Get("/", &response)
	if err != nil {
		return nil, err, queryInfo
	}

	return response.Chains, nil, queryInfo
}
