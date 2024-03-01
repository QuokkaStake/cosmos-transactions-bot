package types

import (
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
)

type Message interface {
	Type() string
	GetAdditionalData(fetcher DataFetcher, subscriptionName string)
	GetValues() event.EventValues
	GetRawMessages() []*codecTypes.Any
	AddParsedMessage(message Message)
	SetParsedMessages(messages []Message)
	GetParsedMessages() []Message
}
