package messages

import (
	configTypes "main/pkg/config/types"
	"main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"
	"main/pkg/utils"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
)

type MsgSend struct {
	From   configTypes.Link
	To     configTypes.Link
	Amount types.Amounts
}

func ParseMsgSend(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosBankTypes.MsgSend
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgSend{
		From:   chain.GetWalletLink(parsedMessage.FromAddress),
		To:     chain.GetWalletLink(parsedMessage.ToAddress),
		Amount: utils.Map(parsedMessage.Amount, types.AmountFrom),
	}, nil
}

func (m MsgSend) Type() string {
	return "/cosmos.bank.v1beta1.MsgSend"
}

func (m *MsgSend) GetAdditionalData(fetcher data_fetcher.DataFetcher) {
	price, found := fetcher.GetPrice()
	if !found {
		return
	}

	for _, amount := range m.Amount {
		if amount.Denom != fetcher.Chain.BaseDenom {
			continue
		}

		amount.AddUSDPrice(fetcher.Chain.DisplayDenom, fetcher.Chain.DenomCoefficient, price)
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
		event.From(cosmosBankTypes.EventTypeTransfer, cosmosBankTypes.AttributeKeySpender, m.From.Value),
		event.From(cosmosBankTypes.EventTypeTransfer, cosmosBankTypes.AttributeKeyRecipient, m.To.Value),
		event.From(cosmosBankTypes.EventTypeCoinSpent, cosmosBankTypes.AttributeKeySpender, m.From.Value),
		event.From(cosmosBankTypes.EventTypeCoinReceived, cosmosBankTypes.AttributeKeyReceiver, m.To.Value),
		event.From(cosmosBankTypes.EventTypeTransfer, cosmosTypes.AttributeKeyAmount, m.Amount.String()),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, m.From.Value),
	}
}
