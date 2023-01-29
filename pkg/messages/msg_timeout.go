package messages

import (
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	ibcTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcChannelTypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/gogo/protobuf/proto"
)

type MsgTimeout struct {
	Signer configTypes.Link
	Packet *Packet
}

func ParseMsgTimeout(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage ibcChannelTypes.MsgTimeout
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	packet, err := ParsePacket(parsedMessage.Packet, chain)
	if err != nil {
		return nil, err
	}

	var packetData ibcTypes.FungibleTokenPacketData
	if err := ibcTypes.ModuleCdc.UnmarshalJSON(parsedMessage.Packet.Data, &packetData); err != nil {
		return nil, err
	}

	return &MsgTimeout{
		Signer: chain.GetWalletLink(parsedMessage.Signer),
		Packet: packet,
	}, nil
}

func (m MsgTimeout) Type() string {
	return "/ibc.core.channel.v1.MsgTimeout"
}

func (m *MsgTimeout) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	if alias := fetcher.AliasManager.Get(fetcher.Chain.Name, m.Signer.Value); alias != "" {
		m.Signer.Title = alias
	}

	m.Packet.GetAdditionalData(fetcher)
}

func (m *MsgTimeout) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
	}
}

func (m *MsgTimeout) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgTimeout) AddParsedMessage(message types.Message) {
}

func (m *MsgTimeout) GetParsedMessages() []types.Message {
	return []types.Message{}
}
