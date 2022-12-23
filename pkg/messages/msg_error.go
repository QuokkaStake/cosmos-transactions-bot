package messages

import (
	"main/pkg/data_fetcher"
)

type MsgError struct {
	Error error
}

func (m MsgError) Type() string {
	return "MsgError"
}

func (m *MsgError) GetAdditionalData(fetcher data_fetcher.DataFetcher) {

}

func (m *MsgError) GetValues() map[string]string {
	return map[string]string{}
}
