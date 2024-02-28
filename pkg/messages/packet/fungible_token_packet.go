package packet

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/types/amount"
	"main/pkg/types/event"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcChannelTypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	ibcTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

type FungibleTokenPacket struct {
	Token    *amount.Amount
	Sender   *configTypes.Link
	Receiver *configTypes.Link

	SrcPort    string
	SrcChannel string
	DstPort    string
	DstChannel string

	Chain *configTypes.Chain
}

func ParseFungibleTokenPacket(
	packetData ibcTypes.FungibleTokenPacketData,
	packet ibcChannelTypes.Packet,
	chain *configTypes.Chain,
) types.Message {
	return &FungibleTokenPacket{
		Token:      amount.AmountFromString(packetData.Amount, packetData.Denom),
		Sender:     &configTypes.Link{Value: packetData.Sender},
		Receiver:   chain.GetWalletLink(packetData.Receiver),
		SrcPort:    packet.SourcePort,
		SrcChannel: packet.SourceChannel,
		DstPort:    packet.DestinationPort,
		DstChannel: packet.DestinationChannel,
		Chain:      chain,
	}
}

func (p FungibleTokenPacket) Type() string {
	return "FungibleTokenPacket"
}

func (p *FungibleTokenPacket) GetAdditionalData(fetcher types.DataFetcher) {
	p.FetchRemoteChainData(fetcher)
	fetcher.PopulateMultichainWallet(p.Chain, p.DstChannel, p.DstPort, p.Sender)

	if alias := fetcher.GetAliasManager().Get(p.Chain.Name, p.Receiver.Value); alias != "" {
		p.Receiver.Title = alias
	}
}

func (p *FungibleTokenPacket) FetchRemoteChainData(fetcher types.DataFetcher) {
	// p.Sender is always someone from the remote chain, so we need to fetch the data
	// from cross-chain.
	// p.Receiver is on native chain, so we can use p.Chain to generate links
	// and denoms and prices.

	trace := ibcTypes.ParseDenomTrace(p.Token.Denom.String())
	p.Token.Denom = amount.Denom(trace.BaseDenom)
	p.Token.BaseDenom = amount.Denom(trace.BaseDenom)

	if !trace.IsNativeDenom() {
		fetcher.PopulateAmount(p.Chain.ChainID, p.Token)
		return
	}

	originalChainID, fetched := fetcher.GetIbcRemoteChainID(p.Chain.ChainID, p.DstChannel, p.DstPort)

	if !fetched {
		return
	}

	if chain, found := fetcher.FindChainById(originalChainID); found {
		fetcher.PopulateAmount(chain.ChainID, p.Token)
	} else {
		fetcher.PopulateAmount(originalChainID, p.Token)
	}
}
func (p *FungibleTokenPacket) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(ibcTypes.EventTypePacket, cosmosBankTypes.AttributeKeyReceiver, p.Receiver.Value),
		event.From(ibcTypes.EventTypePacket, cosmosTypes.AttributeKeyAmount, p.Token.Value.String()),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, p.Sender.Value),
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
