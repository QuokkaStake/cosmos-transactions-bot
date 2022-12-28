package messages

import (
	"main/pkg/data_fetcher"
	"main/pkg/types/event"
)

type MsgUnsupportedMessage struct {
	MsgType string
}

func (m MsgUnsupportedMessage) Type() string {
	return "MsgUnsupportedMessage"
}

func (m *MsgUnsupportedMessage) GetAdditionalData(fetcher data_fetcher.DataFetcher) {
}

func (m *MsgUnsupportedMessage) GetValues() event.EventValues {
	return []event.EventValue{}
}
