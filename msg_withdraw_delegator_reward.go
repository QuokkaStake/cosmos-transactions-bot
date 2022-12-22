package main

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/gogo/protobuf/proto"
)

type MsgWithdrawDelegatorReward struct {
	DelegatorAddress Link
	ValidatorAddress Link
	Amount           []cosmosTypes.Coin
}

func ParseMsgWithdrawDelegatorReward(data []byte, chain *Chain) (MsgWithdrawDelegatorReward, error) {
	var parsedMessage cosmosDistributionTypes.MsgWithdrawDelegatorReward
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return MsgWithdrawDelegatorReward{}, err
	}

	return MsgWithdrawDelegatorReward{
		DelegatorAddress: chain.GetWalletLink(parsedMessage.DelegatorAddress),
		ValidatorAddress: chain.GetWalletLink(parsedMessage.ValidatorAddress),
	}, nil
}

func (m MsgWithdrawDelegatorReward) Type() string {
	return "MsgWithdrawDelegatorReward"
}
