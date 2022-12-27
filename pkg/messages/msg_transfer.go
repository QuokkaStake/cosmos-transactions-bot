package messages

import (
	ibcTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/gogo/protobuf/proto"
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"
)

type MsgTransfer struct {
	Token    *types.Amount
	Sender   configTypes.Link
	Receiver configTypes.Link
}

func ParseMsgTransfer(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage ibcTypes.MsgTransfer
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgTransfer{
		Token: &types.Amount{
			Value: float64(parsedMessage.Token.Amount.Int64()),
			Denom: parsedMessage.Token.Denom,
		},
		Sender:   chain.GetWalletLink(parsedMessage.Sender),
		Receiver: configTypes.Link{Value: parsedMessage.Receiver},
	}, nil
}

func (m MsgTransfer) Type() string {
	return "MsgTransfer"
}

func (m *MsgTransfer) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	price, found := fetcher.GetPrice()
	if found && m.Token.Denom == fetcher.Chain.BaseDenom {
		m.Token.Value /= float64(fetcher.Chain.DenomCoefficient)
		m.Token.Denom = fetcher.Chain.DisplayDenom
		m.Token.PriceUSD = m.Token.Value * price
	}
}

func (m *MsgTransfer) GetValues() event.EventValues {
	return []event.EventValue{
		{Key: "type", Value: "MsgTransfer"},
		{Key: "sender", Value: m.Sender.Value},
		{Key: "receiver", Value: m.Receiver.Value},
	}
}
