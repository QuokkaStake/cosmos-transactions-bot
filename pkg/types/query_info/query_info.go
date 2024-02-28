package query_info

import "time"

type QueryType string

const (
	QueryTypeRewards                  QueryType = "rewards"
	QueryTypeCommission               QueryType = "commission"
	QueryTypeProposal                 QueryType = "proposal"
	QueryTypeStakingParams            QueryType = "staking_params"
	QueryTypeValidator                QueryType = "validator"
	QueryTypeIbcChannel               QueryType = "ibc_channel"
	QueryTypeIbcConnectionClientState QueryType = "ibc_connection_client_state"
	QueryTypeIbcDenomTrace            QueryType = "ibc_denom_trace"
	QueryTypeChainsList               QueryType = "chains_list"
	QueryTypePrices                   QueryType = "prices"
)

type QueryInfo struct {
	Success bool
	Time    time.Duration
	Node    string
}
