package messages

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/gogo/protobuf/proto"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types/chains"
)

type MsgWithdrawDelegatorReward struct {
	DelegatorAddress chains.Link
	ValidatorAddress chains.Link
	Amount           []cosmosTypes.Coin
}

func ParseMsgWithdrawDelegatorReward(data []byte, chain *chains.Chain) (*MsgWithdrawDelegatorReward, error) {
	var parsedMessage cosmosDistributionTypes.MsgWithdrawDelegatorReward
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgWithdrawDelegatorReward{
		DelegatorAddress: chain.GetWalletLink(parsedMessage.DelegatorAddress),
		ValidatorAddress: chain.GetWalletLink(parsedMessage.ValidatorAddress),
	}, nil
}

func (m MsgWithdrawDelegatorReward) Type() string {
	return "MsgWithdrawDelegatorReward"
}

func (m *MsgWithdrawDelegatorReward) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	validator, found := fetcher.GetValidator(m.ValidatorAddress.Title)
	if !found {
		return
	}

	m.ValidatorAddress.Title = validator.Description.Moniker
}
