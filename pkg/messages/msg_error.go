package messages

import (
	"main/pkg/data_fetcher"
	"main/pkg/types/event"
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
