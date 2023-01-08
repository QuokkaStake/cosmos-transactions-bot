package messages

import (
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosAuthzTypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/gogo/protobuf/proto"
)

type MsgExec struct {
	Grantee     configTypes.Link
	RawMessages []*codecTypes.Any
	Messages    []types.Message
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
	}, nil
}

func (m MsgExec) Type() string {
	return "/cosmos.authz.v1beta1.MsgExec"
}

func (m *MsgExec) GetAdditionalData(fetcher dataFetcher.DataFetcher) {
	if alias := fetcher.AliasManager.Get(fetcher.Chain.Name, m.Grantee.Value); alias != "" {
		m.Grantee.Title = alias
	}

	for _, message := range m.Messages {
		if message != nil {
			message.GetAdditionalData(fetcher)
		}
	}
}

func (m *MsgExec) GetValues() event.EventValues {
	values := []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
	}

	return values
}
