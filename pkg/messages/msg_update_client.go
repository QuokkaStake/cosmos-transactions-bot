package messages

import (
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	ibcClientTypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/gogo/protobuf/proto"
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
	return "/ibc.core.client.v1.MsgUpdateClient"
}

func (m *MsgUpdateClient) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	if alias := fetcher.AliasManager.Get(fetcher.Chain.Name, m.Signer.Value); alias != "" {
		m.Signer.Title = alias
	}
}

func (m *MsgUpdateClient) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(ibcClientTypes.EventTypeUpdateClient, ibcClientTypes.AttributeKeyClientID, m.ClientID),
	}
}

func (m *MsgUpdateClient) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgUpdateClient) AddParsedMessage(message types.Message) {
}

func (m *MsgUpdateClient) GetParsedMessages() []types.Message {
	return []types.Message{}
}
