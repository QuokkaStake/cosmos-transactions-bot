package packet

import (
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	ibcTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
)

type FungibleTokenPacket struct {
	Token    *types.Amount
	Sender   configTypes.Link
	Receiver configTypes.Link
}

func ParseFungibleTokenPacket(
	packetData ibcTypes.FungibleTokenPacketData,
	chain *configTypes.Chain,
) types.Message {
	return &FungibleTokenPacket{
		Token:    types.AmountFromString(packetData.Amount, packetData.Denom),
		Sender:   configTypes.Link{Value: packetData.Sender},
		Receiver: chain.GetWalletLink(packetData.Receiver),
	}
}

func (p FungibleTokenPacket) Type() string {
	return "FungibleTokenPacket"
}

func (p *FungibleTokenPacket) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	trace := ibcTypes.ParseDenomTrace(p.Token.Denom)
	p.Token.Denom = trace.BaseDenom

	price, found := fetcher.GetPrice()
	if found && p.Token.Denom == fetcher.Chain.BaseDenom {
		p.Token.AddUSDPrice(fetcher.Chain.DisplayDenom, fetcher.Chain.DenomCoefficient, price)
	}

	if alias := fetcher.AliasManager.Get(fetcher.Chain.Name, p.Receiver.Value); alias != "" {
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
