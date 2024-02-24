package responses

import "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

type IbcDenomTraceResponse struct {
	DenomTrace types.DenomTrace `json:"denom_trace"`
}
