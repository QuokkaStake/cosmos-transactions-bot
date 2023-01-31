package messages

import (
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/gogo/protobuf/proto"
)

type MsgSetWithdrawAddress struct {
	DelegatorAddress configTypes.Link
	WithdrawAddress  configTypes.Link
}

func ParseMsgSetWithdrawAddress(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosDistributionTypes.MsgSetWithdrawAddress
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgSetWithdrawAddress{
		DelegatorAddress: chain.GetWalletLink(parsedMessage.DelegatorAddress),
		WithdrawAddress:  chain.GetValidatorLink(parsedMessage.WithdrawAddress),
	}, nil
}

func (m MsgSetWithdrawAddress) Type() string {
	return "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress"
}

func (m *MsgSetWithdrawAddress) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	if alias := fetcher.AliasManager.Get(fetcher.Chain.Name, m.DelegatorAddress.Value); alias != "" {
		m.DelegatorAddress.Title = alias
	}

	if alias := fetcher.AliasManager.Get(fetcher.Chain.Name, m.WithdrawAddress.Value); alias != "" {
		m.WithdrawAddress.Title = alias
	}
}

func (m *MsgSetWithdrawAddress) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
	}
}

func (m *MsgSetWithdrawAddress) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgSetWithdrawAddress) AddParsedMessage(message types.Message) {
}

func (m *MsgSetWithdrawAddress) GetParsedMessages() []types.Message {
	return []types.Message{}
}
