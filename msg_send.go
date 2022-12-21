package main

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
)

type MsgSend struct {
	From   string
	To     string
	Amount []cosmosTypes.Coin
}

func ParseMsgSend(data []byte) (MsgSend, error) {
	var parsedMessage cosmosBankTypes.MsgSend
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return MsgSend{}, err
	}

	return MsgSend{
		From:   parsedMessage.FromAddress,
		To:     parsedMessage.ToAddress,
		Amount: parsedMessage.Amount,
	}, nil
}

func (m MsgSend) Type() string {
	return "MsgSend"
}
