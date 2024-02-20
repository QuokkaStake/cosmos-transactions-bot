package messages

import (
	"main/pkg/types"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
)

type MsgUnparsedMessage struct {
	MsgType string
	Error   error
}

func (m MsgUnparsedMessage) Type() string {
	return "MsgUnparsedMessage"
}

func (m *MsgUnparsedMessage) GetAdditionalData(fetcher types.DataFetcher) {
}

func (m *MsgUnparsedMessage) GetValues() event.EventValues {
	return []event.EventValue{}
}

func (m *MsgUnparsedMessage) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgUnparsedMessage) AddParsedMessage(message types.Message) {
}

func (m *MsgUnparsedMessage) SetParsedMessages(messages []types.Message) {
}

func (m *MsgUnparsedMessage) GetParsedMessages() []types.Message {
	return []types.Message{}
}
