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

func GetQueryTypes() []QueryType {
	return []QueryType{
		QueryTypeRewards,
		QueryTypeCommission,
		QueryTypeProposal,
		QueryTypeStakingParams,
		QueryTypeValidator,
	}
}

type QueryInfo struct {
	Success bool
	Time    time.Duration
	Node    string
}
