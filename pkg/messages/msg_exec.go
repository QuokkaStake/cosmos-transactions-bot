package messages

import (
	"fmt"
	configTypes "main/pkg/config/types"
	"main/pkg/types"
	"main/pkg/types/event"
	"strconv"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosAuthzTypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/gogo/protobuf/proto"
)

type MsgExec struct {
	Grantee     *configTypes.Link
	RawMessages []*codecTypes.Any
	Messages    []types.Message

	Chain *configTypes.Chain
}

func ParseMsgExec(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosAuthzTypes.MsgExec
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgExec{
		Grantee:     chain.GetWalletLink(parsedMessage.Grantee),
		Messages:    make([]types.Message, 0),
		RawMessages: parsedMessage.Msgs,
		Chain:       chain,
	}, nil
}

func (m MsgExec) Type() string {
	return "/cosmos.authz.v1beta1.MsgExec"
}

func (m *MsgExec) GetAdditionalData(fetcher types.DataFetcher, subscriptionName string) {
	fetcher.PopulateWalletAlias(m.Chain, m.Grantee, subscriptionName)

	for _, message := range m.Messages {
		if message != nil {
			message.GetAdditionalData(fetcher, subscriptionName)
		}
	}
}

func (m *MsgExec) GetValues() event.EventValues {
	values := []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, m.Grantee.Value),
	}

	return values
}

func (m *MsgExec) GetMessagesLabel() string {
	if len(m.Messages) == len(m.RawMessages) {
		return strconv.Itoa(len(m.Messages))
	}

	return fmt.Sprintf("%d, %d skipped", len(m.RawMessages), len(m.RawMessages)-len(m.Messages))
}

func (m *MsgExec) GetRawMessages() []*codecTypes.Any {
	return m.RawMessages
}

func (m *MsgExec) AddParsedMessage(message types.Message) {
	m.Messages = append(m.Messages, message)
}

func (m *MsgExec) SetParsedMessages(messages []types.Message) {
	m.Messages = messages
}

func (m *MsgExec) GetParsedMessages() []types.Message {
	return m.Messages
}
