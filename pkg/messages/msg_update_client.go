package messages

import (
	ibcClientTypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/gogo/protobuf/proto"
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"
)

type MsgUpdateClient struct {
	ClientID string
	Signer   configTypes.Link
}

func ParseMsgUpdateClient(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage ibcClientTypes.MsgUpdateClient
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgUpdateClient{
		ClientID: parsedMessage.ClientId,
		Signer:   chain.GetWalletLink(parsedMessage.Signer),
	}, nil
}

func (m MsgUpdateClient) Type() string {
	return "MsgUpdateClient"
}

func (m *MsgUpdateClient) GetAdditionalData(fetcher dataFetcher.DataFetcher) {

}

func (m *MsgUpdateClient) GetValues() event.EventValues {
	return []event.EventValue{
		{Key: "type", Value: "MsgUpdateClient"},
		{Key: "signer", Value: m.Signer.Value},
		{Key: "client_id", Value: m.ClientID},
	}
}