package types

import (
	"main/pkg/data_fetcher"
	"main/pkg/types/chains"
)

type Tx struct {
	Hash     chains.Link
	Memo     string
	Height   chains.Link
	Messages []Message
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
