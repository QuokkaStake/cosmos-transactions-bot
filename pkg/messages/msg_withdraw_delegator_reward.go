package messages

import (
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/gogo/protobuf/proto"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/chains"
	"main/pkg/utils"
)

type MsgWithdrawDelegatorReward struct {
	DelegatorAddress chains.Link
	ValidatorAddress chains.Link
	Height           int64
	Amount           []*types.Amount
}

func ParseMsgWithdrawDelegatorReward(data []byte, chain *chains.Chain, height int64) (*MsgWithdrawDelegatorReward, error) {
	var parsedMessage cosmosDistributionTypes.MsgWithdrawDelegatorReward
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgWithdrawDelegatorReward{
		DelegatorAddress: chain.GetWalletLink(parsedMessage.DelegatorAddress),
		ValidatorAddress: chain.GetValidatorLink(parsedMessage.ValidatorAddress),
		Height:           height,
	}, nil
}

func (m MsgWithdrawDelegatorReward) Type() string {
	return "MsgWithdrawDelegatorReward"
}

func (m *MsgWithdrawDelegatorReward) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	rewards, found := fetcher.GetRewardsAtBlock(
		m.DelegatorAddress.Value,
		m.ValidatorAddress.Value,
		m.Height,
	)
	if found {
		m.Amount = make([]*types.Amount, len(rewards))

		for index, reward := range rewards {
			m.Amount[index] = &types.Amount{
				Value: utils.StrToFloat64(reward.Amount),
				Denom: reward.Denom,
			}
		}

		price, priceFound := fetcher.GetPrice()
		if priceFound {
			for _, amount := range m.Amount {
				if amount.Denom != fetcher.Chain.BaseDenom {
					continue
				}

				amount.Value /= float64(fetcher.Chain.DenomCoefficient)
				amount.Denom = fetcher.Chain.DisplayDenom
				amount.PriceUSD = amount.Value * price
			}
		}
	}

	if validator, found := fetcher.GetValidator(m.ValidatorAddress.Value); found {
		m.ValidatorAddress.Title = validator.Description.Moniker
	}
}

func (m *MsgWithdrawDelegatorReward) GetValues() map[string]string {
	return map[string]string{
		"type":      "MsgWithdrawDelegatorReward",
		"delegator": m.DelegatorAddress.Value,
		"validator": m.ValidatorAddress.Value,
	}
}
