package packet

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/types/amount"
	"main/pkg/types/event"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	ibcTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

type FungibleTokenPacket struct {
	Token    *amount.Amount
	Sender   configTypes.Link
	Receiver configTypes.Link
}

func ParseFungibleTokenPacket(
	packetData ibcTypes.FungibleTokenPacketData,
	chain *configTypes.Chain,
) types.Message {
	return &FungibleTokenPacket{
		Token:    amount.AmountFromString(packetData.Amount, packetData.Denom),
		Sender:   configTypes.Link{Value: packetData.Sender},
		Receiver: chain.GetWalletLink(packetData.Receiver),
	}
}

func (p FungibleTokenPacket) Type() string {
	return "FungibleTokenPacket"
}

func (p *FungibleTokenPacket) GetAdditionalData(fetcher types.DataFetcher) {
	trace := ibcTypes.ParseDenomTrace(p.Token.Denom)
	p.Token.Denom = trace.BaseDenom

	fetcher.PopulateAmount(p.Token)

	if alias := fetcher.GetAliasManager().Get(fetcher.GetChain().Name, p.Receiver.Value); alias != "" {
		p.Receiver.Title = alias
	}
}

func (p *FungibleTokenPacket) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(ibcTypes.EventTypePacket, cosmosBankTypes.AttributeKeyReceiver, p.Receiver.Value),
		event.From(ibcTypes.EventTypePacket, cosmosTypes.AttributeKeyAmount, p.Token.Value.String()),
	}
}

func (p *FungibleTokenPacket) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (p *FungibleTokenPacket) AddParsedMessage(message types.Message) {
}

func (p *FungibleTokenPacket) SetParsedMessages(messages []types.Message) {
}

func (p *FungibleTokenPacket) GetParsedMessages() []types.Message {
	return []types.Message{}
}
