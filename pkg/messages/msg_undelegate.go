package messages

import (
	"main/pkg/types/amount"
	"time"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"

	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
)

type MsgUndelegate struct {
	DelegatorAddress     configTypes.Link
	ValidatorAddress     configTypes.Link
	UndelegateFinishTime time.Time
	Amount               *amount.Amount
}

func ParseMsgUndelegate(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosStakingTypes.MsgUndelegate
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgUndelegate{
		DelegatorAddress: chain.GetWalletLink(parsedMessage.DelegatorAddress),
		ValidatorAddress: chain.GetValidatorLink(parsedMessage.ValidatorAddress),
		Amount:           amount.AmountFrom(parsedMessage.Amount),
	}, nil
}

func (m MsgUndelegate) Type() string {
	return "/cosmos.staking.v1beta1.MsgUndelegate"
}

func (m *MsgUndelegate) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	if validator, found := fetcher.GetValidator(m.ValidatorAddress.Value); found {
		m.ValidatorAddress.Title = validator.Description.Moniker
	}

	if stakingParams, found := fetcher.GetStakingParams(); found {
		m.UndelegateFinishTime = time.Now().Add(stakingParams.UnbondingTime.Duration)
	}

	fetcher.PopulateAmount(m.Amount)

	if alias := fetcher.AliasManager.Get(fetcher.Chain.Name, m.DelegatorAddress.Value); alias != "" {
		m.DelegatorAddress.Title = alias
	}
}

func (m *MsgUndelegate) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosStakingTypes.EventTypeUnbond, cosmosStakingTypes.AttributeKeyValidator, m.ValidatorAddress.Value),
		event.From(cosmosStakingTypes.EventTypeUnbond, cosmosStakingTypes.AttributeKeyDelegator, m.DelegatorAddress.Value),
		event.From(cosmosStakingTypes.EventTypeUnbond, cosmosTypes.AttributeKeyAmount, m.Amount.String()),
	}
}

func (m *MsgUndelegate) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgUndelegate) AddParsedMessage(message types.Message) {
}

func (m *MsgUndelegate) SetParsedMessages(messages []types.Message) {
}

func (m *MsgUndelegate) GetParsedMessages() []types.Message {
	return []types.Message{}
}
