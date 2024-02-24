package types

import (
	"main/pkg/config/types"
)

type Report struct {
	Chain             *types.Chain
	Subscription      *types.Subscription
	ChainSubscription *types.ChainSubscription
	Node              string
	Reportable        Reportable
}
