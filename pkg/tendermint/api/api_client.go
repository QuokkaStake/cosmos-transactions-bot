package api

import (
	"fmt"
	"main/pkg/http"
	"main/pkg/types/query_info"
	"strconv"

	"github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	configTypes "main/pkg/config/types"
	"main/pkg/types/responses"

	"github.com/rs/zerolog"
)

type TendermintApiClient struct {
	Logger zerolog.Logger
	Client *http.Client
}

func NewTendermintApiClient(logger *zerolog.Logger, url string, chain *configTypes.Chain) *TendermintApiClient {
	return &TendermintApiClient{
		Logger: logger.With().
			Str("component", "tendermint_api_client").
			Str("chain", chain.Name).
			Logger(),
		Client: http.NewClient(logger, url, chain.Name),
	}
}

func (c *TendermintApiClient) GetValidator(address string) (*responses.Validator, error, query_info.QueryInfo) {
	url := fmt.Sprintf("/cosmos/staking/v1beta1/validators/%s", address)

	var response *responses.ValidatorResponse
	err, queryInfo := c.Client.Get(url, &response)
	if err != nil {
		return nil, err, queryInfo
	}

	return &response.Validator, nil, queryInfo
}

func (c *TendermintApiClient) GetDelegatorsRewardsAtBlock(
	delegator string,
	validator string,
	block int64,
) ([]responses.Reward, error, query_info.QueryInfo) {
	url := fmt.Sprintf(
		"/cosmos/distribution/v1beta1/delegators/%s/rewards/%s",
		delegator,
		validator,
	)

	headers := map[string]string{
		"x-cosmos-block-height": strconv.FormatInt(block, 10),
	}

	var response *responses.RewardsResponse
	err, queryInfo := c.Client.GetWithHeaders(url, &response, headers)
	if err != nil || response == nil {
		return nil, err, queryInfo
	}

	return response.Rewards, nil, queryInfo
}

func (c *TendermintApiClient) GetValidatorCommissionAtBlock(
	validator string,
	block int64,
) ([]responses.Commission, error, query_info.QueryInfo) {
	url := fmt.Sprintf(
		"/cosmos/distribution/v1beta1/validators/%s/commission",
		validator,
	)

	headers := map[string]string{
		"x-cosmos-block-height": strconv.FormatInt(block, 10),
	}

	var response *responses.CommissionResponse
	err, queryInfo := c.Client.GetWithHeaders(url, &response, headers)
	if err != nil || response == nil {
		return nil, err, queryInfo
	}

	return response.Commission.Commission, nil, queryInfo
}

func (c *TendermintApiClient) GetProposal(id string) (*responses.Proposal, error, query_info.QueryInfo) {
	url := fmt.Sprintf("/cosmos/gov/v1beta1/proposals/%s", id)

	var response *responses.ProposalResponse
	err, queryInfo := c.Client.Get(url, &response)
	if err != nil {
		return nil, err, queryInfo
	}

	return &response.Proposal, nil, queryInfo
}

func (c *TendermintApiClient) GetStakingParams() (*responses.StakingParams, error, query_info.QueryInfo) {
	url := fmt.Sprintf("/cosmos/staking/v1beta1/params")

	var response *responses.StakingParamsResponse
	err, queryInfo := c.Client.Get(url, &response)
	if err != nil {
		return nil, err, queryInfo
	}

	return &response.Params, nil, queryInfo
}

func (c *TendermintApiClient) GetIbcChannel(
	channel string,
	port string,
) (*responses.IbcChannel, error, query_info.QueryInfo) {
	url := fmt.Sprintf("/ibc/core/channel/v1/channels/%s/ports/%s", channel, port)

	var response *responses.IbcChannelResponse
	err, queryInfo := c.Client.Get(url, &response)
	if err != nil {
		return nil, err, queryInfo
	}

	return &response.Channel, nil, queryInfo
}

func (c *TendermintApiClient) GetIbcConnectionClientState(
	connectionID string,
) (*responses.IbcIdentifiedClientState, error, query_info.QueryInfo) {
	url := fmt.Sprintf("/ibc/core/connection/v1/connections/%s/client_state", connectionID)

	var response *responses.IbcClientStateResponse
	err, queryInfo := c.Client.Get(url, &response)
	if err != nil {
		return nil, err, queryInfo
	}

	return &response.IdentifiedClientState, nil, queryInfo
}

func (c *TendermintApiClient) GetIbcDenomTrace(
	hash string,
) (*types.DenomTrace, error, query_info.QueryInfo) {
	url := fmt.Sprintf("/ibc/apps/transfer/v1/denom_traces/%s", hash)

	var response *responses.IbcDenomTraceResponse
	err, queryInfo := c.Client.Get(url, &response)
	if err != nil {
		return nil, err, queryInfo
	}

	return &response.DenomTrace, nil, queryInfo
}
