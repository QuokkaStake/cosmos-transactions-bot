package types

import (
	"main/pkg/config/types"
	"main/pkg/data_fetcher"
)

type Tx struct {
	Hash          types.Link
	Memo          string
	Height        types.Link
	MessagesCount int
	Messages      []Message
}

func (tx Tx) GetMessages() []Message {
	return tx.Messages
}

func (tx Tx) Type() string {
	return "Tx"
}

func (tx Tx) GetHash() string {
	return tx.Hash.Title
}

func (tx *Tx) GetAdditionalData(fetcher data_fetcher.DataFetcher) {
	for _, msg := range tx.Messages {
		msg.GetAdditionalData(fetcher)
	}
}