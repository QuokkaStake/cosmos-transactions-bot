package api

import (
	"fmt"
	"main/pkg/http"
	"main/pkg/metrics"
	"main/pkg/types/query_info"
	"strconv"

	"github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	configTypes "main/pkg/config/types"
	"main/pkg/types/responses"

	"github.com/rs/zerolog"
)

type TendermintApiClient struct {
	Logger         zerolog.Logger
	Client         *http.Client
	MetricsManager *metrics.Manager
	ChainName      string
}

func NewTendermintApiClient(
	logger *zerolog.Logger,
	url string,
	chain *configTypes.Chain,
	metricsManager *metrics.Manager,
) *TendermintApiClient {
	return &TendermintApiClient{
		Logger: logger.With().
			Str("component", "tendermint_api_client").
			Str("chain", chain.Name).
			Logger(),
		Client:         http.NewClient(logger, url, chain.Name),
		ChainName:      chain.Name,
		MetricsManager: metricsManager,
	}
}

func (c *TendermintApiClient) GetValidator(address string) (*responses.Validator, error) {
	url := fmt.Sprintf("/cosmos/staking/v1beta1/validators/%s", address)

	var response *responses.ValidatorResponse
	err, queryInfo := c.Client.Get(url, &response)
	c.MetricsManager.LogQuery(c.ChainName, queryInfo, query_info.QueryTypeValidator)

	if err != nil {
		return nil, err
	}

	return &response.Validator, nil
}

func (c *TendermintApiClient) GetDelegatorsRewardsAtBlock(
	delegator string,
	validator string,
	block int64,
) ([]responses.Reward, error) {
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
	c.MetricsManager.LogQuery(c.ChainName, queryInfo, query_info.QueryTypeRewards)

	if err != nil || response == nil {
		return nil, err
	}

	return response.Rewards, nil
}

func (c *TendermintApiClient) GetValidatorCommissionAtBlock(
	validator string,
	block int64,
) ([]responses.Commission, error) {
	url := fmt.Sprintf(
		"/cosmos/distribution/v1beta1/validators/%s/commission",
		validator,
	)

	headers := map[string]string{
		"x-cosmos-block-height": strconv.FormatInt(block, 10),
	}

	var response *responses.CommissionResponse
	err, queryInfo := c.Client.GetWithHeaders(url, &response, headers)
	c.MetricsManager.LogQuery(c.ChainName, queryInfo, query_info.QueryTypeCommission)

	if err != nil || response == nil {
		return nil, err
	}

	return response.Commission.Commission, nil
}

func (c *TendermintApiClient) GetProposal(id string) (*responses.Proposal, error) {
	url := fmt.Sprintf("/cosmos/gov/v1beta1/proposals/%s", id)

	var response *responses.ProposalResponse
	err, queryInfo := c.Client.Get(url, &response)
	c.MetricsManager.LogQuery(c.ChainName, queryInfo, query_info.QueryTypeProposal)

	if err != nil {
		return nil, err
	}

	return &response.Proposal, nil
}

func (c *TendermintApiClient) GetStakingParams() (*responses.StakingParams, error) {
	url := fmt.Sprintf("/cosmos/staking/v1beta1/params")

	var response *responses.StakingParamsResponse
	err, queryInfo := c.Client.Get(url, &response)
	c.MetricsManager.LogQuery(c.ChainName, queryInfo, query_info.QueryTypeStakingParams)

	if err != nil {
		return nil, err
	}

	return &response.Params, nil
}

func (c *TendermintApiClient) GetIbcChannel(
	channel string,
	port string,
) (*responses.IbcChannel, error) {
	url := fmt.Sprintf("/ibc/core/channel/v1/channels/%s/ports/%s", channel, port)

	var response *responses.IbcChannelResponse
	err, queryInfo := c.Client.Get(url, &response)
	c.MetricsManager.LogQuery(c.ChainName, queryInfo, query_info.QueryTypeIbcChannel)

	if err != nil {
		return nil, err
	}

	return &response.Channel, nil
}

func (c *TendermintApiClient) GetIbcConnectionClientState(
	connectionID string,
) (*responses.IbcIdentifiedClientState, error) {
	url := fmt.Sprintf("/ibc/core/connection/v1/connections/%s/client_state", connectionID)

	var response *responses.IbcClientStateResponse
	err, queryInfo := c.Client.Get(url, &response)
	c.MetricsManager.LogQuery(c.ChainName, queryInfo, query_info.QueryTypeIbcConnectionClientState)

	if err != nil {
		return nil, err
	}

	return &response.IdentifiedClientState, nil
}

func (c *TendermintApiClient) GetIbcDenomTrace(
	hash string,
) (*types.DenomTrace, error) {
	url := fmt.Sprintf("/ibc/apps/transfer/v1/denom_traces/%s", hash)

	var response *responses.IbcDenomTraceResponse
	err, queryInfo := c.Client.Get(url, &response)
	c.MetricsManager.LogQuery(c.ChainName, queryInfo, query_info.QueryTypeIbcDenomTrace)

	if err != nil {
		return nil, err
	}

	return &response.DenomTrace, nil
}
