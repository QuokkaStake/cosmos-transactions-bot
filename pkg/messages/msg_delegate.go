package messages

import (
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
)

type MsgDelegate struct {
	DelegatorAddress configTypes.Link
	ValidatorAddress configTypes.Link
	Amount           *types.Amount
}

func ParseMsgDelegate(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
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
