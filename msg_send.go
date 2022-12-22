package main

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
)

type MsgSend struct {
	From   Link
	To     Link
	Amount []cosmosTypes.Coin
}

func ParseMsgSend(data []byte, chain *Chain) (MsgSend, error) {
	var parsedMessage cosmosBankTypes.MsgSend
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return MsgSend{}, err
	}

	return MsgSend{
		From:   chain.GetWalletLink(parsedMessage.FromAddress),
		To:     chain.GetWalletLink(parsedMessage.ToAddress),
		Amount: parsedMessage.Amount,
	}, nil
}

func (m MsgSend) Type() string {
	return "MsgSend"
}
