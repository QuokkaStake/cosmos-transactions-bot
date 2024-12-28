package messages

import (
	"main/pkg/types"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
)

type MsgNotExistingMessage struct{}

func (m *MsgNotExistingMessage) Type() string {
	return "MsgNotExistingMessage"
}

func (m *MsgNotExistingMessage) GetAdditionalData(fetcher types.DataFetcher, subscriptionName string) {
}

func (m *MsgNotExistingMessage) GetValues() event.EventValues {
	return []event.EventValue{}
}

func (m *MsgNotExistingMessage) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgNotExistingMessage) AddParsedMessage(message types.Message) {
}

func (m *MsgNotExistingMessage) SetParsedMessages(messages []types.Message) {
}

func (m *MsgNotExistingMessage) GetParsedMessages() []types.Message {
	return []types.Message{}
}
