package types

import (
	"fmt"
	"strconv"

	"main/pkg/config/types"
	"main/pkg/data_fetcher"
)

type Tx struct {
	Hash          types.Link
	Memo          string
	Height        types.Link
	MessagesCount int
	Code          uint32
	Log           string

	Messages []Message
}

func (tx Tx) GetMessages() []Message {
	return tx.Messages
}

func (tx Tx) Type() string {
	return "Tx"
}

func (tx Tx) GetHash() string {
	return tx.Hash.Value
}

func (tx *Tx) GetAdditionalData(fetcher data_fetcher.DataFetcher) {
	for _, msg := range tx.Messages {
		msg.GetAdditionalData(fetcher)
	}
}

func (tx *Tx) GetMessagesLabel() string {
	if len(tx.Messages) == tx.MessagesCount {
		return strconv.Itoa(tx.MessagesCount)
	}

	return fmt.Sprintf("%d, %d skipped", tx.MessagesCount, tx.MessagesCount-len(tx.Messages))
}
