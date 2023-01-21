package messages

import (
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/gogo/protobuf/proto"
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
			m.Amount[index] = types.AmountFromString(reward.Amount, reward.Denom)
		}

		price, priceFound := fetcher.GetPrice()
		if priceFound {
			for _, amount := range m.Amount {
				if amount.Denom != fetcher.Chain.BaseDenom {
					continue
				}

				amount.AddUSDPrice(fetcher.Chain.DisplayDenom, fetcher.Chain.DenomCoefficient, price)
			}
		}
	}

	if validator, found := fetcher.GetValidator(m.ValidatorAddress.Value); found {
		m.ValidatorAddress.Title = validator.Description.Moniker
	}

	if alias := fetcher.AliasManager.Get(fetcher.Chain.Name, m.DelegatorAddress.Value); alias != "" {
		m.DelegatorAddress.Title = alias
	}
}

func (m *MsgWithdrawDelegatorReward) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosDistributionTypes.EventTypeWithdrawRewards, cosmosDistributionTypes.AttributeKeyValidator, m.ValidatorAddress.Value),
	}
}

func (m *MsgWithdrawDelegatorReward) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgWithdrawDelegatorReward) AddParsedMessage(message types.Message) {
}

func (m *MsgWithdrawDelegatorReward) GetParsedMessages() []types.Message {
	return []types.Message{}
}
