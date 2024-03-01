package messages

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosAuthzTypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/gogo/protobuf/proto"
)

type MsgRevoke struct {
	Granter       *configTypes.Link
	Grantee       *configTypes.Link
	MsgType       string
	Authorization Authorization

	Chain *configTypes.Chain
}

func ParseMsgRevoke(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosAuthzTypes.MsgRevoke
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgRevoke{
		Grantee: chain.GetWalletLink(parsedMessage.Grantee),
		Granter: chain.GetWalletLink(parsedMessage.Granter),
		MsgType: parsedMessage.MsgTypeUrl,
		Chain:   chain,
	}, nil
}

func (m MsgRevoke) Type() string {
	return "/cosmos.authz.v1beta1.MsgRevoke"
}

func (m *MsgRevoke) GetAdditionalData(fetcher types.DataFetcher, subscriptionName string) {
	fetcher.PopulateWalletAlias(m.Chain, m.Grantee, subscriptionName)
	fetcher.PopulateWalletAlias(m.Chain, m.Granter, subscriptionName)
}

func (m *MsgRevoke) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, m.Grantee.Value),
	}
}

func (m *MsgRevoke) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgRevoke) AddParsedMessage(message types.Message) {
}

func (m *MsgRevoke) SetParsedMessages(messages []types.Message) {
}

func (m *MsgRevoke) GetParsedMessages() []types.Message {
	return []types.Message{}
}
