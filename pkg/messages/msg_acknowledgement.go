package messages

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/types/amount"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	ibcTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcChannelTypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/gogo/protobuf/proto"
)

type MsgAcknowledgement struct {
	Token    *amount.Amount
	Sender   *configTypes.Link
	Receiver *configTypes.Link
	Signer   *configTypes.Link

	Chain *configTypes.Chain
}

func ParseMsgAcknowledgement(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage ibcChannelTypes.MsgAcknowledgement
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	var packetData ibcTypes.FungibleTokenPacketData
	if err := ibcTypes.ModuleCdc.UnmarshalJSON(parsedMessage.Packet.Data, &packetData); err != nil {
		return nil, err
	}

	return &MsgAcknowledgement{
		Token:    amount.AmountFromString(packetData.Amount, packetData.Denom),
		Sender:   chain.GetWalletLink(packetData.Sender),
		Receiver: &configTypes.Link{Value: packetData.Receiver},
		Signer:   chain.GetWalletLink(parsedMessage.Signer),
		Chain:    chain,
	}, nil
}

func (m MsgAcknowledgement) Type() string {
	return "/ibc.core.channel.v1.MsgAcknowledgement"
}

func (m *MsgAcknowledgement) GetAdditionalData(fetcher types.DataFetcher) {
	fetcher.PopulateAmount(m.Chain.ChainID, m.Token)
	if alias := fetcher.GetAliasManager().Get(m.Chain.Name, m.Sender.Value); alias != "" {
		m.Sender.Title = alias
	}

	if alias := fetcher.GetAliasManager().Get(m.Chain.Name, m.Signer.Value); alias != "" {
		m.Signer.Title = alias
	}
}

func (m *MsgAcknowledgement) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, m.Signer.Value),
	}
}

func (m *MsgAcknowledgement) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgAcknowledgement) AddParsedMessage(message types.Message) {
}

func (m *MsgAcknowledgement) SetParsedMessages(messages []types.Message) {
}

func (m *MsgAcknowledgement) GetParsedMessages() []types.Message {
	return []types.Message{}
}
