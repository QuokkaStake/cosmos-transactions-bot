package types

import (
	"main/pkg/config/types"
)

type Report struct {
	Chain      types.Chain
	Node       string
	Reportable Reportable
}
