package main

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
)

type MsgSend struct {
	From   Link
	To     Link
	Amount []*Amount
}

func ParseMsgSend(data []byte, chain *Chain) (*MsgSend, error) {
	var parsedMessage cosmosBankTypes.MsgSend
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgSend{
		From: chain.GetWalletLink(parsedMessage.FromAddress),
		To:   chain.GetWalletLink(parsedMessage.ToAddress),
		Amount: Map(parsedMessage.Amount, func(coin cosmosTypes.Coin) *Amount {
			return &Amount{
				Value: float64(coin.Amount.Int64()),
				Denom: coin.Denom,
			}
		}),
	}, nil
}

func (m MsgSend) Type() string {
	return "MsgSend"
}

func (m *MsgSend) GetAdditionalData(fetcher DataFetcher) {
	price, found := fetcher.GetPrice()
	if !found {
		return
	}

	for _, amount := range m.Amount {
		if amount.Denom != fetcher.Chain.BaseDenom {
			continue
		}

		amount.Value /= float64(fetcher.Chain.DenomCoefficient)
		amount.Denom = fetcher.Chain.DisplayDenom
		amount.PriceUSD = amount.Value * price
	}
}
