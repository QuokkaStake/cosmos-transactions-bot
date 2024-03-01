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

type MsgBeginRedelegate struct {
	DelegatorAddress    *configTypes.Link
	ValidatorSrcAddress *configTypes.Link
	ValidatorDstAddress *configTypes.Link
	Amount              *amount.Amount

	Chain *configTypes.Chain
}

func ParseMsgBeginRedelegate(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosStakingTypes.MsgBeginRedelegate
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgBeginRedelegate{
		DelegatorAddress:    chain.GetWalletLink(parsedMessage.DelegatorAddress),
		ValidatorSrcAddress: chain.GetValidatorLink(parsedMessage.ValidatorSrcAddress),
		ValidatorDstAddress: chain.GetValidatorLink(parsedMessage.ValidatorDstAddress),
		Amount:              amount.AmountFrom(parsedMessage.Amount),
		Chain:               chain,
	}, nil
}

func (m MsgBeginRedelegate) Type() string {
	return "/cosmos.staking.v1beta1.MsgBeginRedelegate"
}

func (m *MsgBeginRedelegate) GetAdditionalData(fetcher types.DataFetcher, subscriptionName string) {
	fetcher.PopulateValidator(m.Chain, m.ValidatorSrcAddress)
	fetcher.PopulateValidator(m.Chain, m.ValidatorDstAddress)

	fetcher.PopulateAmount(m.Chain.ChainID, m.Amount)
	fetcher.PopulateWalletAlias(m.Chain, m.DelegatorAddress, subscriptionName)
}

func (m *MsgBeginRedelegate) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, m.DelegatorAddress.Value),
		event.From(cosmosStakingTypes.EventTypeRedelegate, cosmosStakingTypes.AttributeKeySrcValidator, m.ValidatorSrcAddress.Value),
		event.From(cosmosStakingTypes.EventTypeRedelegate, cosmosStakingTypes.AttributeKeyDstValidator, m.ValidatorDstAddress.Value),
		event.From(cosmosStakingTypes.EventTypeRedelegate, cosmosStakingTypes.AttributeKeyDelegator, m.DelegatorAddress.Value),
		event.From(cosmosStakingTypes.EventTypeRedelegate, cosmosTypes.AttributeKeyAmount, m.Amount.String()),
	}
}

func (m *MsgBeginRedelegate) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgBeginRedelegate) AddParsedMessage(message types.Message) {
}

func (m *MsgBeginRedelegate) SetParsedMessages(messages []types.Message) {
}

func (m *MsgBeginRedelegate) GetParsedMessages() []types.Message {
	return []types.Message{}
}
