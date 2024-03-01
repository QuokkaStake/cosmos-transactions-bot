package messages

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/types/amount"
	"main/pkg/types/event"
	"main/pkg/utils"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
)

type MsgSend struct {
	From   *configTypes.Link
	To     *configTypes.Link
	Amount amount.Amounts

	Chain *configTypes.Chain
}

func ParseMsgSend(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosBankTypes.MsgSend
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgSend{
		From:   chain.GetWalletLink(parsedMessage.FromAddress),
		To:     chain.GetWalletLink(parsedMessage.ToAddress),
		Amount: utils.Map(parsedMessage.Amount, amount.AmountFrom),
		Chain:  chain,
	}, nil
}

func (m MsgSend) Type() string {
	return "/cosmos.bank.v1beta1.MsgSend"
}

func (m *MsgSend) GetAdditionalData(fetcher types.DataFetcher, subscriptionName string) {
	fetcher.PopulateAmounts(m.Chain.ChainID, m.Amount)

	fetcher.PopulateWalletAlias(m.Chain, m.From, subscriptionName)
	fetcher.PopulateWalletAlias(m.Chain, m.To, subscriptionName)
}

func (m *MsgSend) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, m.From.Value),
		event.From(cosmosBankTypes.EventTypeTransfer, cosmosTypes.AttributeKeySender, m.From.Value),
		event.From(cosmosBankTypes.EventTypeTransfer, cosmosBankTypes.AttributeKeyRecipient, m.To.Value),
		event.From(cosmosBankTypes.EventTypeCoinSpent, cosmosBankTypes.AttributeKeySpender, m.From.Value),
		event.From(cosmosBankTypes.EventTypeCoinReceived, cosmosBankTypes.AttributeKeyReceiver, m.To.Value),
		event.From(cosmosBankTypes.EventTypeTransfer, cosmosTypes.AttributeKeyAmount, m.Amount.String()),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, m.From.Value),
	}
}

func (m *MsgSend) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgSend) AddParsedMessage(message types.Message) {
}

func (m *MsgSend) SetParsedMessages(messages []types.Message) {
}

func (m *MsgSend) GetParsedMessages() []types.Message {
	return []types.Message{}
}
