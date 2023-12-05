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

type MsgWithdrawValidatorCommission struct {
	ValidatorAddress configTypes.Link
	Height           int64
	Amount           []*amount.Amount
}

func ParseMsgWithdrawValidatorCommission(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosDistributionTypes.MsgWithdrawValidatorCommission
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgWithdrawValidatorCommission{
		ValidatorAddress: chain.GetValidatorLink(parsedMessage.ValidatorAddress),
		Height:           height,
	}, nil
}

func (m MsgWithdrawValidatorCommission) Type() string {
	return "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"
}

func (m *MsgWithdrawValidatorCommission) GetAdditionalData(fetcher types.DataFetcher) {
	rewards, found := fetcher.GetCommissionAtBlock(
		m.ValidatorAddress.Value,
		m.Height,
	)
	if found {
		m.Amount = make([]*amount.Amount, len(rewards))

		for index, reward := range rewards {
			m.Amount[index] = amount.AmountFromString(reward.Amount, reward.Denom)
		}

		fetcher.PopulateAmounts(m.Amount)
	}

	if validator, found := fetcher.GetValidator(m.ValidatorAddress.Value); found {
		m.ValidatorAddress.Title = validator.Description.Moniker
	}
}

func (m *MsgWithdrawValidatorCommission) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosDistributionTypes.EventTypeWithdrawCommission, cosmosDistributionTypes.AttributeKeyValidator, m.ValidatorAddress.Value),
	}
}

func (m *MsgWithdrawValidatorCommission) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgWithdrawValidatorCommission) AddParsedMessage(message types.Message) {
}

func (m *MsgWithdrawValidatorCommission) SetParsedMessages(messages []types.Message) {
}

func (m *MsgWithdrawValidatorCommission) GetParsedMessages() []types.Message {
	return []types.Message{}
}
