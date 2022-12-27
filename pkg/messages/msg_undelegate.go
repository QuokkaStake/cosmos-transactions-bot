package messages

import (
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"
	"time"
)

type MsgUndelegate struct {
	DelegatorAddress     configTypes.Link
	ValidatorAddress     configTypes.Link
	UndelegateFinishTime time.Time
	Amount               *types.Amount
}

func ParseMsgUndelegate(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosStakingTypes.MsgUndelegate
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgUndelegate{
		DelegatorAddress: chain.GetWalletLink(parsedMessage.DelegatorAddress),
		ValidatorAddress: chain.GetValidatorLink(parsedMessage.ValidatorAddress),
		Amount: &types.Amount{
			Value: float64(parsedMessage.Amount.Amount.Int64()),
			Denom: parsedMessage.Amount.Denom,
		},
	}, nil
}

func (m MsgUndelegate) Type() string {
	return "MsgUndelegate"
}

func (m *MsgUndelegate) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	if validator, found := fetcher.GetValidator(m.ValidatorAddress.Value); found {
		m.ValidatorAddress.Title = validator.Description.Moniker
	}

	if stakingParams, found := fetcher.GetStakingParams(); found {
		m.UndelegateFinishTime = time.Now().Add(stakingParams.UnbondingTime.Duration)
	}

	price, found := fetcher.GetPrice()
	if found && m.Amount.Denom == fetcher.Chain.BaseDenom {
		m.Amount.Value /= float64(fetcher.Chain.DenomCoefficient)
		m.Amount.Denom = fetcher.Chain.DisplayDenom
		m.Amount.PriceUSD = m.Amount.Value * price
	}
}

func (m *MsgUndelegate) GetValues() event.EventValues {
	return []event.EventValue{
		{Key: "type", Value: "MsgUndelegate"},
		{Key: "delegator_address", Value: m.DelegatorAddress.Value},
		{Key: "validator_address", Value: m.ValidatorAddress.Value},
	}
}
