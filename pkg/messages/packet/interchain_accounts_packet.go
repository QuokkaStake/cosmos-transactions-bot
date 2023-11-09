package packet

import (
	"fmt"
	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/types/event"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/std"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	icaTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
)

type InterchainAccountsPacket struct {
	PacketType      string
	Memo            string
	TxMessagesCount int
	TxRawMessages   []*codecTypes.Any
	TxMessages      []types.Message
}

func ParseInterchainAccountsPacket(
	packetData icaTypes.InterchainAccountPacketData,
	chain *configTypes.Chain,
) (types.Message, error) {
	cdc := codec.NewLegacyAmino()
	interfaceRegistry := codecTypes.NewInterfaceRegistry()
	codec := codec.NewProtoCodec(interfaceRegistry)

	std.RegisterLegacyAminoCodec(cdc)
	std.RegisterInterfaces(interfaceRegistry)

	var cosmosTx icaTypes.CosmosTx
	if err := codec.Unmarshal(packetData.Data, &cosmosTx); err != nil {
		return nil, err
	}

	return &InterchainAccountsPacket{
		PacketType:      packetData.Type.String(),
		Memo:            packetData.Memo,
		TxMessagesCount: len(cosmosTx.Messages),
		TxRawMessages:   cosmosTx.Messages,
		TxMessages:      make([]types.Message, 0),
	}, nil
}

func (p InterchainAccountsPacket) Type() string {
	return "InterchainAccountsPacket"
}

func (p *InterchainAccountsPacket) GetAdditionalData(fetcher types.DataFetcher) {
	for _, message := range p.TxMessages {
		message.GetAdditionalData(fetcher)
	}
}

func (p *InterchainAccountsPacket) GetValues() event.EventValues {
	var values []event.EventValue

	for _, message := range p.TxMessages {
		values = append(values, message.GetValues()...)
	}

	return values
}

func (p *InterchainAccountsPacket) GetRawMessages() []*codecTypes.Any {
	return p.TxRawMessages
}

func (p *InterchainAccountsPacket) AddParsedMessage(message types.Message) {
	p.TxMessages = append(p.TxMessages, message)
}

func (p *InterchainAccountsPacket) SetParsedMessages(messages []types.Message) {
	p.TxMessages = messages
}

func (p *InterchainAccountsPacket) GetParsedMessages() []types.Message {
	return p.TxMessages
}

func (p *InterchainAccountsPacket) GetMessagesLabel() string {
	if len(p.TxMessages) == len(p.TxRawMessages) {
		return strconv.Itoa(len(p.TxMessages))
	}

	return fmt.Sprintf("%d, %d skipped", len(p.TxRawMessages), len(p.TxRawMessages)-len(p.TxMessages))
}
