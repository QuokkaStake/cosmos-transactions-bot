package messages

import (
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	ibcTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcChannelTypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
)

type Packet struct {
	Token    *types.Amount
	Sender   configTypes.Link
	Receiver configTypes.Link
}

func ParsePacket(packet ibcChannelTypes.Packet, chain *configTypes.Chain) (*Packet, error) {
	var packetData ibcTypes.FungibleTokenPacketData
	if err := ibcTypes.ModuleCdc.UnmarshalJSON(packet.Data, &packetData); err != nil {
		return nil, err
	}

	return &Packet{
		Token:    types.AmountFromString(packetData.Amount, packetData.Denom),
		Sender:   configTypes.Link{Value: packetData.Sender},
		Receiver: chain.GetWalletLink(packetData.Receiver),
	}, nil
}

func (p Packet) Type() string {
	return "Packet"
}

func (p *Packet) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
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

func (p *Packet) GetValues() event.EventValues {
	return []event.EventValue{}
}

func (p *Packet) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (p *Packet) AddParsedMessage(message types.Message) {
}

func (p *Packet) GetParsedMessages() []types.Message {
	return []types.Message{}
}
