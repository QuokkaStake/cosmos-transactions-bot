package types

import (
	"main/pkg/config/types"
)

type MessageParser func([]byte, *types.Chain, int64) (Message, error)
