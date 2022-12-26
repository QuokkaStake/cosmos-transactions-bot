package responses

import (
	"encoding/json"
	"fmt"
	"time"
)

type ValidatorResponse struct {
	Validator Validator `json:"validator"`
}

type Validator struct {
	OperatorAddress   string               `json:"operator_address"`
	ConsensusPubkey   ConsensusPubkey      `json:"consensus_pubkey"`
	Jailed            bool                 `json:"jailed"`
	Status            string               `json:"status"`
	Tokens            string               `json:"tokens"`
	DelegatorShares   string               `json:"delegator_shares"`
	Description       ValidatorDescription `json:"description"`
	UnbondingHeight   string               `json:"unbonding_height"`
	UnbondingTime     time.Time            `json:"unbonding_time"`
	Commission        ValidatorCommission  `json:"commission"`
	MinSelfDelegation string               `json:"min_self_delegation"`
}

type ConsensusPubkey struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type ValidatorDescription struct {
	Moniker         string `json:"moniker"`
	Identity        string `json:"identity"`
	Website         string `json:"website"`
	SecurityContact string `json:"security_contact"`
	Details         string `json:"details"`
}

type ValidatorCommission struct {
	CommissionRates ValidatorCommissionRates `json:"commission_rates"`
	UpdateTime      time.Time                `json:"update_time"`
}

type ValidatorCommissionRates struct {
	Rate          string `json:"rate"`
	MaxRate       string `json:"max_rate"`
	MaxChangeRate string `json:"max_change_rate"`
}

type RewardsResponse struct {
	Rewards []Reward `json:"rewards"`
}

type Reward struct {
	Amount string `json:"amount"`
	Denom  string `json:"denom"`
}

type ProposalResponse struct {
	Proposal Proposal `json:"proposal"`
}

type Proposal struct {
	ProposalID    string          `json:"proposal_id"`
	Content       ProposalContent `json:"content"`
	Status        string          `json:"status"`
	VotingEndTime time.Time       `json:"voting_end_time"`
}

type ProposalContent struct {
	Type        string `json:"@type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type StakingParamsResponse struct {
	Params StakingParams `json:"params"`
}

type StakingParams struct {
	UnbondingTime Duration `json:"unbonding_time"`
}

// Golang cannot properly deserialize string into time.Duration, that's why this workaround.
// Cheers to https://biscuit.ninja/posts/go-unmarshalling-json-into-time-duration/
type Duration struct {
	time.Duration
}

func (duration *Duration) UnmarshalJSON(b []byte) error {
	var unmarshalledJson interface{}

	err := json.Unmarshal(b, &unmarshalledJson)
	if err != nil {
		return err
	}

	switch value := unmarshalledJson.(type) {
	case string:
		duration.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid duration: %#v", unmarshalledJson)
	}

	return nil
}
