package messages

import (
	"main/pkg/data_fetcher"
)

type MsgUnsupportedMessage struct {
	MsgType string
}

func (m MsgUnsupportedMessage) Type() string {
	return "MsgUnsupportedMessage"
}

func (m *MsgUnsupportedMessage) GetAdditionalData(fetcher data_fetcher.DataFetcher) {

}

func (m *MsgUnsupportedMessage) GetValues() map[string]string {
	return map[string]string{}
}
