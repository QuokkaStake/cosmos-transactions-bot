package messages

import (
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
	types2 "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
)

type MsgDelegate struct {
	DelegatorAddress types2.Link
	ValidatorAddress types2.Link
	Amount           *types.Amount
}

func ParseMsgDelegate(data []byte, chain *types2.Chain) (*MsgDelegate, error) {
	var parsedMessage cosmosStakingTypes.MsgDelegate
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgDelegate{
		DelegatorAddress: chain.GetWalletLink(parsedMessage.DelegatorAddress),
		ValidatorAddress: chain.GetValidatorLink(parsedMessage.ValidatorAddress),
		Amount: &types.Amount{
			Value: float64(parsedMessage.Amount.Amount.Int64()),
			Denom: parsedMessage.Amount.Denom,
		},
	}, nil
}

func (m MsgDelegate) Type() string {
	return "MsgDelegate"
}

func (m *MsgDelegate) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	validator, found := fetcher.GetValidator(m.ValidatorAddress.Value)
	if found {
		m.ValidatorAddress.Title = validator.Description.Moniker
	}

	price, found := fetcher.GetPrice()
	if found && m.Amount.Denom == fetcher.Chain.BaseDenom {
		m.Amount.Value /= float64(fetcher.Chain.DenomCoefficient)
		m.Amount.Denom = fetcher.Chain.DisplayDenom
		m.Amount.PriceUSD = m.Amount.Value * price
	}
}

func (m *MsgDelegate) GetValues() map[string]string {
	return map[string]string{
		"type":              "MsgDelegate",
		"delegator_address": m.DelegatorAddress.Value,
		"validator_address": m.ValidatorAddress.Value,
	}
}
