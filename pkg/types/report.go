package types

import (
	"main/pkg/config/types"
)

type Report struct {
	Chain        types.Chain
	Subscription types.Subscription
	Node         string
	Reportable   Reportable
}
