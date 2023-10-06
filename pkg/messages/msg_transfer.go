package messages

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/types/amount"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	ibcTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/gogo/protobuf/proto"
)

type MsgTransfer struct {
	Token    *amount.Amount
	Sender   configTypes.Link
	Receiver configTypes.Link
}

func ParseMsgTransfer(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage ibcTypes.MsgTransfer
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgTransfer{
		Token:    amount.AmountFrom(parsedMessage.Token),
		Sender:   chain.GetWalletLink(parsedMessage.Sender),
		Receiver: configTypes.Link{Value: parsedMessage.Receiver},
	}, nil
}

func (m MsgTransfer) Type() string {
	return "/ibc.applications.transfer.v1.MsgTransfer"
}

func (m *MsgTransfer) GetAdditionalData(fetcher types.DataFetcher) {
	fetcher.PopulateAmount(m.Token)

	if alias := fetcher.GetAliasManager().Get(fetcher.GetChain().Name, m.Sender.Value); alias != "" {
		m.Sender.Title = alias
	}
}

func (m *MsgTransfer) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
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
