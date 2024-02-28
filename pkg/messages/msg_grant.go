package messages

import (
	"main/pkg/types"
	"main/pkg/types/amount"
	"time"

	configTypes "main/pkg/config/types"
	"main/pkg/types/event"
	"main/pkg/utils"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosAuthzTypes "github.com/cosmos/cosmos-sdk/x/authz"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gogo/protobuf/proto"
)

type Authorization interface{}

type StakeAuthorization struct {
	MaxTokens         *amount.Amount
	AuthorizationType string
	Validators        []configTypes.Link
}

type MsgGrant struct {
	Granter       *configTypes.Link
	Grantee       *configTypes.Link
	GrantType     string
	Expiration    *time.Time
	Authorization Authorization

	Chain *configTypes.Chain
}

func ParseStakeAuthorization(authorization *codecTypes.Any, chain *configTypes.Chain) (Authorization, error) {
	var parsedAuthorization cosmosStakingTypes.StakeAuthorization
	if err := proto.Unmarshal(authorization.Value, &parsedAuthorization); err != nil {
		return nil, err
	}
	var maxTokens *amount.Amount
	if parsedAuthorization.MaxTokens != nil {
		maxTokens = amount.AmountFrom(*parsedAuthorization.MaxTokens)
	}

	var validators []configTypes.Link
	authorizationType := "UNSPECIFIED"

	if allowList := parsedAuthorization.GetAllowList(); allowList != nil {
		validators = utils.Map(allowList.Address, func(address string) configTypes.Link {
			return chain.GetValidatorLink(address)
		})
		authorizationType = "ALLOWLIST"
	} else if denyList := parsedAuthorization.GetDenyList(); denyList != nil {
		validators = utils.Map(denyList.Address, func(address string) configTypes.Link {
			return chain.GetValidatorLink(address)
		})
		authorizationType = "DENYLIST"
	}

	generatedAuthorization := StakeAuthorization{
		MaxTokens:         maxTokens,
		Validators:        validators,
		AuthorizationType: authorizationType,
	}

	return generatedAuthorization, nil
}

func ParseMsgGrant(data []byte, chain *configTypes.Chain, height int64) (types.Message, error) {
	var parsedMessage cosmosAuthzTypes.MsgGrant
	if err := proto.Unmarshal(data, &parsedMessage); err != nil {
		return nil, err
	}

	var authorization Authorization

	if parsedMessage.Grant.Authorization.TypeUrl == "/cosmos.staking.v1beta1.StakeAuthorization" {
		if value, err := ParseStakeAuthorization(parsedMessage.Grant.Authorization, chain); err != nil {
			return nil, err
		} else {
			authorization = value
		}
	}

	return &MsgGrant{
		Grantee:       chain.GetWalletLink(parsedMessage.Grantee),
		Granter:       chain.GetWalletLink(parsedMessage.Granter),
		GrantType:     parsedMessage.Grant.Authorization.TypeUrl,
		Expiration:    parsedMessage.Grant.Expiration,
		Authorization: authorization,
		Chain:         chain,
	}, nil
}

func (m MsgGrant) Type() string {
	return "/cosmos.authz.v1beta1.MsgGrant"
}

func (m *MsgGrant) GetAdditionalData(fetcher types.DataFetcher) {
	if alias := fetcher.GetAliasManager().Get(m.Chain.Name, m.Grantee.Value); alias != "" {
		m.Grantee.Title = alias
	}

	if alias := fetcher.GetAliasManager().Get(m.Chain.Name, m.Granter.Value); alias != "" {
		m.Granter.Title = alias
	}
}

func (m *MsgGrant) GetValues() event.EventValues {
	return []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, m.Type()),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, m.Grantee.Value),
	}
}

func (m *MsgGrant) GetRawMessages() []*codecTypes.Any {
	return []*codecTypes.Any{}
}

func (m *MsgGrant) AddParsedMessage(message types.Message) {
}

func (m *MsgGrant) SetParsedMessages(messages []types.Message) {
}

func (m *MsgGrant) GetParsedMessages() []types.Message {
	return []types.Message{}
}
