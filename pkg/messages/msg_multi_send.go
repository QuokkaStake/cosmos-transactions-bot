package messages

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
	configTypes "main/pkg/config/types"
	"main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"
	"main/pkg/utils"
)

type MultiSendEntry struct {
	Address configTypes.Link
	Amount  []*types.Amount
}

type MsgMultiSend struct {
	Inputs  []MultiSendEntry
	Outputs []MultiSendEntry
}

func ParseMsgMultiSend(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosBankTypes.MsgMultiSend
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgMultiSend{
		Inputs: utils.Map(parsedMessage.Inputs, func(input cosmosBankTypes.Input) MultiSendEntry {
			return MultiSendEntry{
				Address: chain.GetWalletLink(input.Address),
				Amount: utils.Map(input.Coins, func(coin cosmosTypes.Coin) *types.Amount {
					return &types.Amount{
						Value: float64(coin.Amount.Int64()),
						Denom: coin.Denom,
					}
				}),
			}
		}),
		Outputs: utils.Map(parsedMessage.Outputs, func(output cosmosBankTypes.Output) MultiSendEntry {
			return MultiSendEntry{
				Address: chain.GetWalletLink(output.Address),
				Amount: utils.Map(output.Coins, func(coin cosmosTypes.Coin) *types.Amount {
					return &types.Amount{
						Value: float64(coin.Amount.Int64()),
						Denom: coin.Denom,
					}
				}),
			}
		}),
	}, nil
}

func (m MsgMultiSend) Type() string {
	return "/cosmos.bank.v1beta1.MsgMultiSend"
}

func (m *MsgMultiSend) GetAdditionalData(fetcher data_fetcher.DataFetcher) {
	price, found := fetcher.GetPrice()
	if !found {
		return
	}

	for _, input := range m.Inputs {
		for _, amount := range input.Amount {
			if amount.Denom != fetcher.Chain.BaseDenom {
				continue
			}

			amount.Value /= float64(fetcher.Chain.DenomCoefficient)
			amount.Denom = fetcher.Chain.DisplayDenom
			amount.PriceUSD = amount.Value * price
		}
	}

	for _, output := range m.Outputs {
		for _, amount := range output.Amount {
			if amount.Denom != fetcher.Chain.BaseDenom {
				continue
			}

			amount.Value /= float64(fetcher.Chain.DenomCoefficient)
			amount.Denom = fetcher.Chain.DisplayDenom
			amount.PriceUSD = amount.Value * price
		}
	}
}

func (m *MsgMultiSend) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
	}
}
