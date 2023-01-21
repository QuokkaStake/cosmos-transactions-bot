package messages

import (
	"main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
)

type MsgError struct {
	Error error
}

func (m MsgError) Type() string {
	return "MsgError"
}

func (m *MsgError) GetAdditionalData(fetcher data_fetcher.DataFetcher) {
}

func (m *MsgError) GetValues() event.EventValues {
	return []event.EventValue{}
}

func (m *MsgError) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgError) AddParsedMessage(message types.Message) {
}

func (m *MsgError) GetParsedMessages() []types.Message {
	return []types.Message{}
}
