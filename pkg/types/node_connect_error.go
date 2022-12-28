package types

import (
	"main/pkg/data_fetcher"
)

type NodeConnectError struct {
	Error error
	Chain string
	URL   string
}

func (e NodeConnectError) GetMessages() []Message {
	return []Message{}
}

func (e NodeConnectError) Type() string {
	return "NodeConnectError"
}

func (e NodeConnectError) GetHash() string {
	return "NodeConnectError"
}

func (e *NodeConnectError) GetAdditionalData(fetcher data_fetcher.DataFetcher) {

}
