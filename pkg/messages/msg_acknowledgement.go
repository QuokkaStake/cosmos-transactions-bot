package messages

import (
	ibcTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcChannelTypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/gogo/protobuf/proto"
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"
	"main/pkg/utils"
)

type MsgAcknowledgement struct {
	Token    *types.Amount
	Sender   configTypes.Link
	Receiver configTypes.Link
	Signer   configTypes.Link
}

func ParseMsgAcknowledgement(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage ibcChannelTypes.MsgAcknowledgement
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	var packetData ibcTypes.FungibleTokenPacketData
	if err := ibcTypes.ModuleCdc.UnmarshalJSON(parsedMessage.Packet.Data, &packetData); err != nil {
		return nil, err
	}

	return &MsgAcknowledgement{
		Token: &types.Amount{
			Value: utils.StrToFloat64(packetData.Amount),
			Denom: packetData.Denom,
		},
		Sender:   chain.GetWalletLink(packetData.Sender),
		Receiver: configTypes.Link{Value: packetData.Receiver},
		Signer:   chain.GetWalletLink(parsedMessage.Signer),
	}, nil
}

func (m MsgAcknowledgement) Type() string {
	return "MsgAcknowledgement"
}

func (m *MsgAcknowledgement) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	price, found := fetcher.GetPrice()
	if found && m.Token.Denom == fetcher.Chain.BaseDenom {
		m.Token.Value /= float64(fetcher.Chain.DenomCoefficient)
		m.Token.Denom = fetcher.Chain.DisplayDenom
		m.Token.PriceUSD = m.Token.Value * price
	}
}

func (m *MsgAcknowledgement) GetValues() event.EventValues {
	return []event.EventValue{
		{Key: "type", Value: "MsgAcknowledgement"},
		{Key: "sender", Value: m.Sender.Value},
		{Key: "receiver", Value: m.Receiver.Value},
		{Key: "signer", Value: m.Signer.Value},
	}
}
