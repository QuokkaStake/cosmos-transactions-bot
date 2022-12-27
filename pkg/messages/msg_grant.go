package messages

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosAuthzTypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/gogo/protobuf/proto"
	configTypes "main/pkg/config/types"
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types"
	"main/pkg/types/event"
	"time"
)

type Authorization interface{}

type StakeAuthorization struct {
	MaxTokens         *types.Amount
	AuthorizationType string
}

type MsgGrant struct {
	Granter       configTypes.Link
	Grantee       configTypes.Link
	GrantType     string
	Expiration    *time.Time
	Authorization Authorization
}

func ParseMsgGrant(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosAuthzTypes.MsgGrant
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	return &MsgGrant{
		Grantee:    chain.GetWalletLink(parsedMessage.Grantee),
		Granter:    chain.GetValidatorLink(parsedMessage.Granter),
		GrantType:  parsedMessage.Grant.Authorization.TypeUrl,
		Expiration: parsedMessage.Grant.Expiration,
	}, nil
}

func (m MsgGrant) Type() string {
	return "MsgGrant"
}

func (m *MsgGrant) GetAdditionalData(fetcher dataFetcher.DataFetcher) {

}

func (m *MsgGrant) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/cosmos.authz.v1beta1.MsgGrant"),
	}
}
