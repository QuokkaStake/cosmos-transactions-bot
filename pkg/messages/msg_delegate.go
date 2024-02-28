package messages

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/types/amount"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
)

type MsgDelegate struct {
	DelegatorAddress *configTypes.Link
	ValidatorAddress configTypes.Link
	Amount           *amount.Amount

	Chain *configTypes.Chain
}

func ParseMsgDelegate(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosStakingTypes.MsgDelegate
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgDelegate{
		DelegatorAddress: chain.GetWalletLink(parsedMessage.DelegatorAddress),
		ValidatorAddress: chain.GetValidatorLink(parsedMessage.ValidatorAddress),
		Amount:           amount.AmountFrom(parsedMessage.Amount),
		Chain:            chain,
	}, nil
}

func (m MsgDelegate) Type() string {
	return "/cosmos.staking.v1beta1.MsgDelegate"
}

func (m *MsgDelegate) GetAdditionalData(fetcher types.DataFetcher) {
	validator, found := fetcher.GetValidator(m.Chain, m.ValidatorAddress.Value)
	if found {
		m.ValidatorAddress.Title = validator.Description.Moniker
	}

	fetcher.PopulateAmount(m.Chain.ChainID, m.Amount)

	if alias := fetcher.GetAliasManager().Get(m.Chain.Name, m.DelegatorAddress.Value); alias != "" {
		m.DelegatorAddress.Title = alias
	}
}

func (m *MsgDelegate) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, m.DelegatorAddress.Value),
		event.From(cosmosStakingTypes.EventTypeDelegate, cosmosStakingTypes.AttributeKeyValidator, m.ValidatorAddress.Value),
		event.From(cosmosStakingTypes.EventTypeDelegate, cosmosStakingTypes.AttributeKeyDelegator, m.DelegatorAddress.Value),
		event.From(cosmosStakingTypes.EventTypeDelegate, cosmosTypes.AttributeKeyAmount, m.Amount.String()),
	}
}

func (m *MsgDelegate) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgDelegate) AddParsedMessage(message types.Message) {
}

func (m *MsgDelegate) SetParsedMessages(messages []types.Message) {
}

func (m *MsgDelegate) GetParsedMessages() []types.Message {
	return []types.Message{}
}
