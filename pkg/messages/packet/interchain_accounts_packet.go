package packet

import (
	"fmt"
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"
	"strconv"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	icaTypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
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
	simappConfig := simapp.MakeTestEncodingConfig()

	var cosmosTx icaTypes.CosmosTx
	if err := simappConfig.Codec.Unmarshal(packetData.Data, &cosmosTx); err != nil {
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

func (p *InterchainAccountsPacket) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
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

func (p *InterchainAccountsPacket) GetParsedMessages() []types.Message {
	return p.TxMessages
}

func (p *InterchainAccountsPacket) GetMessagesLabel() string {
	if len(p.TxMessages) == len(p.TxRawMessages) {
		return strconv.Itoa(len(p.TxMessages))
	}

	return fmt.Sprintf("%d, %d skipped", len(p.TxRawMessages), len(p.TxRawMessages)-len(p.TxMessages))
}
