package types

import (
	"github.com/google/uuid"
)

type TxError struct {
	Error error
}

func (txError TxError) GetMessages() []Message {
	return []Message{}
}

func (txError TxError) Type() string {
	return "TxError"
}

func (txError TxError) GetHash() string {
	return uuid.NewString()
}

func (txError *TxError) GetAdditionalData(fetcher DataFetcher) {
}
