package query_info

import "time"

type QueryType string

const (
	QueryTypeRewards       QueryType = "rewards"
	QueryTypeCommission    QueryType = "commission"
	QueryTypeProposal      QueryType = "proposal"
	QueryTypeStakingParams QueryType = "staking_params"
	QueryTypeValidator     QueryType = "validator"
)

type QueryInfo struct {
	Success bool
	Time    time.Duration
	Node    string
}
