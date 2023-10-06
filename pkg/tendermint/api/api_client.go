package api

import (
	"encoding/json"
	"fmt"
	"main/pkg/types/query_info"
	"net/http"
	"strconv"
	"time"

	configTypes "main/pkg/config/types"
	"main/pkg/types/responses"

	"github.com/rs/zerolog"
)

type TendermintApiClient struct {
	Logger  zerolog.Logger
	URL     string
	Timeout time.Duration
}

func NewTendermintApiClient(logger *zerolog.Logger, url string, chain *configTypes.Chain) *TendermintApiClient {
	return &TendermintApiClient{
		Logger: logger.With().
			Str("component", "tendermint_api_client").
			Str("url", url).
			Str("chain", chain.Name).
			Logger(),
		URL:     url,
		Timeout: 60 * time.Second,
	}
}

func (c *TendermintApiClient) GetValidator(address string) (*responses.Validator, error, query_info.QueryInfo) {
	url := fmt.Sprintf("/cosmos/staking/v1beta1/validators/%s", address)

	var response *responses.ValidatorResponse
	err, queryInfo := c.Get(url, &response)
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
	err, queryInfo := c.GetWithHeaders(url, &response, headers)
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
	err, queryInfo := c.GetWithHeaders(url, &response, headers)
	if err != nil || response == nil {
		return nil, err, queryInfo
	}

	return response.Commission.Commission, nil, queryInfo
}

func (c *TendermintApiClient) GetProposal(id string) (*responses.Proposal, error, query_info.QueryInfo) {
	url := fmt.Sprintf("/cosmos/gov/v1beta1/proposals/%s", id)

	var response *responses.ProposalResponse
	err, queryInfo := c.Get(url, &response)
	if err != nil {
		return nil, err, queryInfo
	}

	return &response.Proposal, nil, queryInfo
}

func (c *TendermintApiClient) GetStakingParams() (*responses.StakingParams, error, query_info.QueryInfo) {
	url := fmt.Sprintf("/cosmos/staking/v1beta1/params")

	var response *responses.StakingParamsResponse
	err, queryInfo := c.Get(url, &response)
	if err != nil {
		return nil, err, queryInfo
	}

	return &response.Params, nil, queryInfo
}

func (c *TendermintApiClient) Get(url string, target interface{}) (error, query_info.QueryInfo) {
	return c.GetWithHeaders(url, target, map[string]string{})
}

func (c *TendermintApiClient) GetWithHeaders(
	relativeURL string,
	target interface{},
	headers map[string]string,
) (error, query_info.QueryInfo) {
	url := fmt.Sprintf("%s%s", c.URL, relativeURL)

	client := &http.Client{Timeout: c.Timeout}
	start := time.Now()
	queryInfo := query_info.QueryInfo{
		Success: false,
		Node:    c.URL,
		Time:    0,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err, queryInfo
	}

	req.Header.Set("User-Agent", "cosmos-transactions-bot")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	c.Logger.Trace().Str("url", url).Msg("Doing a query...")

	res, err := client.Do(req)
	queryInfo.Time = time.Since(start)
	if err != nil {
		c.Logger.Warn().Str("url", url).Err(err).Msg("Query failed")
		return err, queryInfo
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		c.Logger.Warn().
			Str("url", url).
			Err(err).
			Int("status", res.StatusCode).
			Msg("Query returned bad HTTP code")
		return fmt.Errorf("bad HTTP code: %d", res.StatusCode), queryInfo
	}

	c.Logger.Debug().Str("url", url).Dur("duration", time.Since(start)).Msg("Query is finished")

	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return err, queryInfo
	}

	queryInfo.Success = true
	return nil, queryInfo
}
