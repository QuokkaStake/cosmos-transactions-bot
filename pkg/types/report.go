package types

import "main/pkg/types/chains"

type Report struct {
	Chain      chains.Chain
	Node       string
	Reportable Reportable
}
