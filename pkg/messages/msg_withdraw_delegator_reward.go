package messages

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/types/amount"
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
	Amount           []*amount.Amount

	Chain *configTypes.Chain
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
		Chain:            chain,
	}, nil
}

func (m MsgWithdrawDelegatorReward) Type() string {
	return "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"
}

func (m *MsgWithdrawDelegatorReward) GetAdditionalData(fetcher types.DataFetcher) {
	rewards, found := fetcher.GetRewardsAtBlock(
		m.Chain,
		m.DelegatorAddress.Value,
		m.ValidatorAddress.Value,
		m.Height,
	)
	if found {
		m.Amount = make([]*amount.Amount, len(rewards))

		for index, reward := range rewards {
			m.Amount[index] = amount.AmountFromString(reward.Amount, reward.Denom)
		}

		fetcher.PopulateAmounts(m.Chain, m.Amount)
	}

	if validator, found := fetcher.GetValidator(m.Chain, m.ValidatorAddress.Value); found {
		m.ValidatorAddress.Title = validator.Description.Moniker
	}

	if alias := fetcher.GetAliasManager().Get(m.Chain.Name, m.DelegatorAddress.Value); alias != "" {
		m.DelegatorAddress.Title = alias
	}
}

func (m *MsgWithdrawDelegatorReward) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosDistributionTypes.EventTypeWithdrawRewards, cosmosDistributionTypes.AttributeKeyValidator, m.ValidatorAddress.Value),
		event.From(cosmosDistributionTypes.EventTypeWithdrawRewards, cosmosDistributionTypes.AttributeKeyDelegator, m.DelegatorAddress.Value),
	}
}

func (m *MsgWithdrawDelegatorReward) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgWithdrawDelegatorReward) AddParsedMessage(message types.Message) {
}

func (m *MsgWithdrawDelegatorReward) SetParsedMessages(messages []types.Message) {
}

func (m *MsgWithdrawDelegatorReward) GetParsedMessages() []types.Message {
	return []types.Message{}
}
