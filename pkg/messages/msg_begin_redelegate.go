package messages

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"
)

type MsgBeginRedelegate struct {
	DelegatorAddress    configTypes.Link
	ValidatorSrcAddress configTypes.Link
	ValidatorDstAddress configTypes.Link
	Amount              *types.Amount
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
		Amount: &types.Amount{
			Value: float64(parsedMessage.Amount.Amount.Int64()),
			Denom: parsedMessage.Amount.Denom,
		},
	}, nil
}

func (m MsgBeginRedelegate) Type() string {
	return "MsgBeginRedelegate"
}

func (m *MsgBeginRedelegate) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	if validator, found := fetcher.GetValidator(m.ValidatorSrcAddress.Value); found {
		m.ValidatorSrcAddress.Title = validator.Description.Moniker
	}
	if validator, found := fetcher.GetValidator(m.ValidatorDstAddress.Value); found {
		m.ValidatorDstAddress.Title = validator.Description.Moniker
	}

	price, found := fetcher.GetPrice()
	if found && m.Amount.Denom == fetcher.Chain.BaseDenom {
		m.Amount.Value /= float64(fetcher.Chain.DenomCoefficient)
		m.Amount.Denom = fetcher.Chain.DisplayDenom
		m.Amount.PriceUSD = m.Amount.Value * price
	}
}

func (m *MsgBeginRedelegate) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/cosmos.staking.v1beta1.MsgBeginRedelegate"),
		{Key: "delegator_address", Value: m.DelegatorAddress.Value},
		{Key: "validator_src_address", Value: m.ValidatorSrcAddress.Value},
		{Key: "validator_dst_address", Value: m.ValidatorDstAddress.Value},
	}
}
