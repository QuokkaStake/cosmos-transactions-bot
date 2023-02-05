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

type MsgTimeout struct {
	Signer configTypes.Link
	Packet types.Message
}

func ParseMsgTimeout(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage ibcChannelTypes.MsgTimeout
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	packet, err := packet.ParsePacket(parsedMessage.Packet, chain)
	if err != nil {
		return nil, err
	} else if packet == nil {
		return nil, nil
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
	values := []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
	}

	values = append(values, m.Packet.GetValues()...)
	return values
}

func (m *MsgTimeout) GetRawMessages() []*codecTypes.Any {
	return m.Packet.GetRawMessages()
}

func (m *MsgTimeout) AddParsedMessage(message types.Message) {
	m.Packet.AddParsedMessage(message)
}

func (m *MsgTimeout) SetParsedMessages(messages []types.Message) {
	m.Packet.SetParsedMessages(messages)
}

func (m *MsgTimeout) GetParsedMessages() []types.Message {
	return m.Packet.GetParsedMessages()
}
