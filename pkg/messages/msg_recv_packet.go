package messages

import (
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/messages/packet"
	"main/pkg/types"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	ibcChannelTypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/gogo/protobuf/proto"
)

type MsgRecvPacket struct {
	Signer configTypes.Link
	Packet types.Message
}

func ParseMsgRecvPacket(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage ibcChannelTypes.MsgRecvPacket
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	parsedPacket, err := packet.ParsePacket(parsedMessage.Packet, chain)
	if err != nil {
		return nil, err
	} else if parsedPacket == nil {
		return nil, nil
	}

	return &MsgRecvPacket{
		Signer: chain.GetWalletLink(parsedMessage.Signer),
		Packet: parsedPacket,
	}, nil
}

func (m MsgRecvPacket) Type() string {
	return "/ibc.core.channel.v1.MsgRecvPacket"
}

func (m *MsgRecvPacket) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	if alias := fetcher.AliasManager.Get(fetcher.Chain.Name, m.Signer.Value); alias != "" {
		m.Signer.Title = alias
	}

	m.Packet.GetAdditionalData(fetcher)
}

func (m *MsgRecvPacket) GetValues() event.EventValues {
	values := []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
	}

	values = append(values, m.Packet.GetValues()...)
	return values
}

func (m *MsgRecvPacket) GetRawMessages() []*codecTypes.Any {
	return m.Packet.GetRawMessages()
}

func (m *MsgRecvPacket) AddParsedMessage(message types.Message) {
	m.Packet.AddParsedMessage(message)
}

func (m *MsgRecvPacket) SetParsedMessages(messages []types.Message) {
	m.Packet.SetParsedMessages(messages)
}

func (m *MsgRecvPacket) GetParsedMessages() []types.Message {
	return m.Packet.GetParsedMessages()
}
