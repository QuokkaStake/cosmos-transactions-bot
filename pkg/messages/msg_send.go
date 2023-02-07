package messages

import (
	configTypes "main/pkg/config/types"
	"main/pkg/data_fetcher"
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
	From   configTypes.Link
	To     configTypes.Link
	Amount amount.Amounts
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
	}, nil
}

func (m MsgSend) Type() string {
	return "/cosmos.bank.v1beta1.MsgSend"
}

func (m *MsgSend) GetAdditionalData(fetcher data_fetcher.DataFetcher) {
	for _, amount := range m.Amount {
		fetcher.PopulateAmount(amount)
	}

	if alias := fetcher.AliasManager.Get(fetcher.Chain.Name, m.From.Value); alias != "" {
		m.From.Title = alias
	}

	if alias := fetcher.AliasManager.Get(fetcher.Chain.Name, m.To.Value); alias != "" {
		m.To.Title = alias
	}
}

func (m *MsgSend) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
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
