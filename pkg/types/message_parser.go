package types

import (
	"main/pkg/types/chains"
)

type MessageParser func([]byte, chains.Chain, int64) (Message, error)
