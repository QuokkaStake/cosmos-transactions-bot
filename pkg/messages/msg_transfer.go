package messages

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/types/amount"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	ibcTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/gogo/protobuf/proto"
)

type MsgTransfer struct {
	Token    *amount.Amount
	Sender   *configTypes.Link
	Receiver *configTypes.Link

	SrcChannel string
	SrcPort    string

	Chain *configTypes.Chain
}

func ParseMsgTransfer(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage ibcTypes.MsgTransfer
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgTransfer{
		Token:      amount.AmountFrom(parsedMessage.Token),
		Sender:     chain.GetWalletLink(parsedMessage.Sender),
		Receiver:   &configTypes.Link{Value: parsedMessage.Receiver},
		SrcChannel: parsedMessage.SourceChannel,
		SrcPort:    parsedMessage.SourcePort,
		Chain:      chain,
	}, nil
}

func (m MsgTransfer) Type() string {
	return "/ibc.applications.transfer.v1.MsgTransfer"
}

func (m *MsgTransfer) GetAdditionalData(fetcher types.DataFetcher, subscriptionName string) {
	m.FetchRemoteChainData(fetcher)
	fetcher.PopulateMultichainWallet(m.Chain, m.SrcChannel, m.SrcPort, m.Receiver, subscriptionName)
	fetcher.PopulateWalletAlias(m.Chain, m.Sender, subscriptionName)
}

func (m *MsgTransfer) FetchRemoteChainData(fetcher types.DataFetcher) {
	// p.Receiver is always someone from the remote chain, so we need to fetch the data
	// from cross-chain.
	// p.Sender is on native chain, so we can use p.Chain to generate links
	// and denoms and prices.

	// If it's an IBC token (like, withdraw on Osmosis) - we need to figure out what
	// the original denom is, to convert it, and also take the remote chain for links
	// generation.
	// If it's a native token - just take the denom from the current chain, but also fetch
	// the remote chain for links generation.
	var trace ibcTypes.DenomTrace
	if m.Token.Denom.IsIbcToken() {
		externalTrace, found := fetcher.GetDenomTrace(m.Chain, m.Token.Denom.String())
		if !found {
			return
		}
		trace = *externalTrace
	} else {
		trace = ibcTypes.ParseDenomTrace(m.Token.Denom.String())
	}

	m.Token.Denom = amount.Denom(trace.BaseDenom)
	m.Token.BaseDenom = amount.Denom(trace.BaseDenom)

	// If it's native - populate denom as it is, taking current chain as the source chain.
	if trace.IsNativeDenom() {
		fetcher.PopulateAmount(m.Chain.ChainID, m.Token)
		return
	}

	// If it's not native - we need the remote chain ID to get the original denoms from,
	// if we can't fetch it - we can't fetch prices, or generate links (if chain is in local
	// config.)
	originalChainID, fetched := fetcher.GetIbcRemoteChainID(m.Chain.ChainID, m.SrcChannel, m.SrcPort)
	if !fetched {
		return
	}

	fetcher.PopulateAmount(originalChainID, m.Token)
}

func (m *MsgTransfer) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, m.Sender.Value),
		event.From(ibcTypes.EventTypeTransfer, ibcTypes.AttributeKeyReceiver, m.Receiver.Value),
		event.From(ibcTypes.EventTypeTransfer, cosmosTypes.AttributeKeySender, m.Sender.Value),
		event.From(ibcTypes.EventTypeTransfer, cosmosTypes.AttributeKeyAmount, m.Token.String()),
	}
}

func (m *MsgTransfer) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgTransfer) AddParsedMessage(message types.Message) {
}

func (m *MsgTransfer) SetParsedMessages(messages []types.Message) {
}

func (m *MsgTransfer) GetParsedMessages() []types.Message {
	return []types.Message{}
}
