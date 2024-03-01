package messages

import (
	"main/pkg/types"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
)

type MsgUnsupportedMessage struct {
	MsgType string
}

func (m MsgUnsupportedMessage) Type() string {
	return "MsgUnsupportedMessage"
}

func (m *MsgUnsupportedMessage) GetAdditionalData(fetcher types.DataFetcher, subscriptionName string) {
}

func (m *MsgUnsupportedMessage) GetValues() event.EventValues {
	return []event.EventValue{}
}

func (m *MsgUnsupportedMessage) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgUnsupportedMessage) AddParsedMessage(message types.Message) {
}

func (m *MsgUnsupportedMessage) SetParsedMessages(messages []types.Message) {
}

func (m *MsgUnsupportedMessage) GetParsedMessages() []types.Message {
	return []types.Message{}
}
