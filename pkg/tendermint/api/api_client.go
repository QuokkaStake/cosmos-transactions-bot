package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"main/pkg/config/types"
	"main/pkg/types/responses"

	"github.com/rs/zerolog"
)

type TendermintApiClient struct {
	Logger  zerolog.Logger
	URL     string
	Timeout time.Duration
}

func NewTendermintApiClient(logger *zerolog.Logger, url string, chain *types.Chain) *TendermintApiClient {
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

func (c *TendermintApiClient) GetValidator(address string) (*responses.Validator, error) {
	url := fmt.Sprintf(
		"%s/cosmos/staking/v1beta1/validators/%s",
		c.URL,
		address,
	)

	var response *responses.ValidatorResponse
	if err := c.Get(url, &response); err != nil {
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
		"%s/cosmos/distribution/v1beta1/delegators/%s/rewards/%s",
		c.URL,
		delegator,
		validator,
	)

	headers := map[string]string{
		"x-cosmos-block-height": strconv.FormatInt(block, 10),
	}

	var response *responses.RewardsResponse
	if err := c.GetWithHeaders(url, &response, headers); err != nil || response == nil {
		return nil, err
	}

	return response.Rewards, nil
}

func (c *TendermintApiClient) GetValidatorCommissionAtBlock(
	validator string,
	block int64,
) ([]responses.Commission, error) {
	url := fmt.Sprintf(
		"%s/cosmos/distribution/v1beta1/validators/%s/commission",
		c.URL,
		validator,
	)

	headers := map[string]string{
		"x-cosmos-block-height": strconv.FormatInt(block, 10),
	}

	var response *responses.CommissionResponse
	if err := c.GetWithHeaders(url, &response, headers); err != nil || response == nil {
		return nil, err
	}

	return response.Commission.Commission, nil
}

func (c *TendermintApiClient) GetProposal(id string) (*responses.Proposal, error) {
	url := fmt.Sprintf(
		"%s/cosmos/gov/v1beta1/proposals/%s",
		c.URL,
		id,
	)

	var response *responses.ProposalResponse
	if err := c.Get(url, &response); err != nil {
		return nil, err
	}

	return &response.Proposal, nil
}

func (c *TendermintApiClient) GetStakingParams() (*responses.StakingParams, error) {
	url := fmt.Sprintf("%s/cosmos/staking/v1beta1/params", c.URL)

	var response *responses.StakingParamsResponse
	if err := c.Get(url, &response); err != nil {
		return nil, err
	}

	return &response.Params, nil
}

func (c *TendermintApiClient) Get(url string, target interface{}) error {
	return c.GetWithHeaders(url, target, map[string]string{})
}

func (c *TendermintApiClient) GetWithHeaders(
	url string,
	target interface{},
	headers map[string]string,
) error {
	client := &http.Client{
		Timeout: c.Timeout,
	}
	start := time.Now()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "cosmos-transactions-bot")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	c.Logger.Trace().Str("url", url).Msg("Doing a query...")

	res, err := client.Do(req)
	if err != nil {
		c.Logger.Warn().Str("url", url).Err(err).Msg("Query failed")
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		c.Logger.Warn().
			Str("url", url).
			Err(err).
			Int("status", res.StatusCode).
			Msg("Query returned bad HTTP code")
		return fmt.Errorf("bad HTTP code: %d", res.StatusCode)
	}

	c.Logger.Debug().Str("url", url).Dur("duration", time.Since(start)).Msg("Query is finished")

	return json.NewDecoder(res.Body).Decode(target)
}
