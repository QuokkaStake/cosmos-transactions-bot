package messages

import (
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	ibcTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcChannelTypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/gogo/protobuf/proto"
)

type MsgRecvPacket struct {
	Token    *types.Amount
	Sender   configTypes.Link
	Receiver configTypes.Link
}

func ParseMsgRecvPacket(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage ibcChannelTypes.MsgRecvPacket
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	var packetData ibcTypes.FungibleTokenPacketData
	if err := ibcTypes.ModuleCdc.UnmarshalJSON(parsedMessage.Packet.Data, &packetData); err != nil {
		return nil, err
	}

	return &MsgRecvPacket{
		Token:    types.AmountFromString(packetData.Amount, packetData.Denom),
		Sender:   configTypes.Link{Value: packetData.Sender},
		Receiver: chain.GetWalletLink(packetData.Receiver),
	}, nil
}

func (m MsgRecvPacket) Type() string {
	return "/ibc.core.channel.v1.MsgRecvPacket"
}

func (m *MsgRecvPacket) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	price, found := fetcher.GetPrice()
	if found && m.Token.Denom == fetcher.Chain.BaseDenom {
		m.Token.AddUSDPrice(fetcher.Chain.DisplayDenom, fetcher.Chain.DenomCoefficient, price)
	}

	if alias := fetcher.AliasManager.Get(fetcher.Chain.Name, m.Receiver.Value); alias != "" {
		m.Receiver.Title = alias
	}
}

func (m *MsgRecvPacket) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
	}
}

func (m *MsgRecvPacket) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgRecvPacket) AddParsedMessage(message types.Message) {
}

func (m *MsgRecvPacket) GetParsedMessages() []types.Message {
	return []types.Message{}
}
