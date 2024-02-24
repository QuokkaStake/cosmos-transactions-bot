package messages

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	ibcClientTypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/gogo/protobuf/proto"
)

type MsgUpdateClient struct {
	ClientID string
	Signer   configTypes.Link

	Chain *configTypes.Chain
}

func ParseMsgUpdateClient(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage ibcClientTypes.MsgUpdateClient
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgUpdateClient{
		ClientID: parsedMessage.ClientId,
		Signer:   chain.GetWalletLink(parsedMessage.Signer),
		Chain:    chain,
	}, nil
}

func (m MsgUpdateClient) Type() string {
	return "/ibc.core.client.v1.MsgUpdateClient"
}

func (m *MsgUpdateClient) GetAdditionalData(fetcher types.DataFetcher) {
	if alias := fetcher.GetAliasManager().Get(m.Chain.Name, m.Signer.Value); alias != "" {
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

func (m *MsgUpdateClient) SetParsedMessages(messages []types.Message) {
}

func (m *MsgUpdateClient) GetParsedMessages() []types.Message {
	return []types.Message{}
}
