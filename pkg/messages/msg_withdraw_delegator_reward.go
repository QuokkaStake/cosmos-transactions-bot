package messages

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/gogo/protobuf/proto"
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"
	"main/pkg/utils"
)

type MsgWithdrawDelegatorReward struct {
	DelegatorAddress configTypes.Link
	ValidatorAddress configTypes.Link
	Height           int64
	Amount           []*types.Amount
}

func ParseMsgWithdrawDelegatorReward(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
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
	return "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"
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

func (m *MsgWithdrawDelegatorReward) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosDistributionTypes.EventTypeWithdrawRewards, cosmosDistributionTypes.AttributeKeyValidator, m.ValidatorAddress.Value),
	}
}
