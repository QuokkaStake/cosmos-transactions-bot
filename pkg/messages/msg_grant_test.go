package messages

import (
	"fmt"
	aliasManagerPkg "main/pkg/alias_manager"
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
	"main/pkg/data_fetcher"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"
	"main/pkg/metrics"
	"main/pkg/types"
	"main/pkg/types/event"
	"testing"

	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	cosmosAuthzTypes "github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMsgGrantParseFail(t *testing.T) {
	t.Parallel()

	parsed, err := ParseMsgGrant([]byte("aaa"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err)
	require.Nil(t, parsed)
}

func TestMsgGrantParse(t *testing.T) {
	t.Parallel()

	msgInternal := &cosmosStakingTypes.StakeAuthorization{}
	msgInternalBytes, err := proto.Marshal(msgInternal)
	require.NoError(t, err)

	msg := &cosmosAuthzTypes.MsgGrant{
		Grantee: "grantee",
		Grant: cosmosAuthzTypes.Grant{
			Authorization: &codecTypes.Any{
				TypeUrl: "/cosmos.staking.v1beta1.StakeAuthorization",
				Value:   msgInternalBytes,
			},
		},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgGrant(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)
}

func TestMsgGrantParseStakeAuthorizationFail(t *testing.T) {
	t.Parallel()

	msg := &cosmosAuthzTypes.MsgGrant{
		Grantee: "grantee",
		Grant: cosmosAuthzTypes.Grant{
			Authorization: &codecTypes.Any{
				TypeUrl: "/cosmos.staking.v1beta1.StakeAuthorization",
				Value:   []byte("string"),
			},
		},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgGrant(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err)
	require.Nil(t, parsed)
}

func TestMsgGrantParseAllowlist(t *testing.T) {
	t.Parallel()

	msgInternal := &cosmosStakingTypes.StakeAuthorization{
		MaxTokens: &cosmosTypes.Coin{Denom: "uatom", Amount: cosmosTypes.NewInt(100)},
		Validators: &cosmosStakingTypes.StakeAuthorization_AllowList{
			AllowList: &cosmosStakingTypes.StakeAuthorization_Validators{
				Address: []string{"validator"},
			},
		},
	}
	msgInternalBytes, err := proto.Marshal(msgInternal)
	require.NoError(t, err)

	msg := &cosmosAuthzTypes.MsgGrant{
		Grantee: "grantee",
		Grant: cosmosAuthzTypes.Grant{
			Authorization: &codecTypes.Any{
				TypeUrl: "/cosmos.staking.v1beta1.StakeAuthorization",
				Value:   msgInternalBytes,
			},
		},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgGrant(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	message, _ := parsed.(*MsgGrant)
	require.NotNil(t, message.Authorization)

	authorization, _ := message.Authorization.(StakeAuthorization)
	require.Equal(t, "ALLOWLIST", authorization.AuthorizationType)
	require.Len(t, authorization.Validators, 1)
	require.Equal(t, "validator", authorization.Validators[0].Value)

	require.NotNil(t, authorization.MaxTokens)
	require.Equal(t, "uatom", authorization.MaxTokens.Denom.String())
	require.Equal(t, "100.00", fmt.Sprintf("%.2f", authorization.MaxTokens.Value))
}

func TestMsgGrantParseDenylist(t *testing.T) {
	t.Parallel()

	msgInternal := &cosmosStakingTypes.StakeAuthorization{
		Validators: &cosmosStakingTypes.StakeAuthorization_DenyList{
			DenyList: &cosmosStakingTypes.StakeAuthorization_Validators{
				Address: []string{"validator"},
			},
		},
	}
	msgInternalBytes, err := proto.Marshal(msgInternal)
	require.NoError(t, err)

	msg := &cosmosAuthzTypes.MsgGrant{
		Grantee: "grantee",
		Grant: cosmosAuthzTypes.Grant{
			Authorization: &codecTypes.Any{
				TypeUrl: "/cosmos.staking.v1beta1.StakeAuthorization",
				Value:   msgInternalBytes,
			},
		},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgGrant(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	message, _ := parsed.(*MsgGrant)
	require.NotNil(t, message.Authorization)

	authorization, _ := message.Authorization.(StakeAuthorization)
	require.Equal(t, "DENYLIST", authorization.AuthorizationType)
	require.Len(t, authorization.Validators, 1)
	require.Equal(t, "validator", authorization.Validators[0].Value)
}

func TestMsgGrantBase(t *testing.T) {
	t.Parallel()

	msgInternal := &cosmosStakingTypes.StakeAuthorization{}
	msgInternalBytes, err := proto.Marshal(msgInternal)
	require.NoError(t, err)

	msg := &cosmosAuthzTypes.MsgGrant{
		Grantee: "grantee",
		Grant: cosmosAuthzTypes.Grant{
			Authorization: &codecTypes.Any{
				TypeUrl: "/cosmos.staking.v1beta1.StakeAuthorization",
				Value:   msgInternalBytes,
			},
		},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgGrant(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/cosmos.authz.v1beta1.MsgGrant", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/cosmos.authz.v1beta1.MsgGrant"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "grantee"),
	}, values)

	parsed.SetParsedMessages([]types.Message{})
	parsed.AddParsedMessage(nil)
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgGrantPopulate(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains:      configTypes.Chains{{Name: "chain"}},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	msgInternal := &cosmosStakingTypes.StakeAuthorization{}
	msgInternalBytes, err := proto.Marshal(msgInternal)
	require.NoError(t, err)

	msg := &cosmosAuthzTypes.MsgGrant{
		Grantee: "grantee",
		Granter: "granter",
		Grant: cosmosAuthzTypes.Grant{
			Authorization: &codecTypes.Any{
				TypeUrl: "/cosmos.staking.v1beta1.StakeAuthorization",
				Value:   msgInternalBytes,
			},
		},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgGrant(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	err = aliasManager.Set("subscription", "chain", "granter", "granter_alias")
	require.NoError(t, err)
	err = aliasManager.Set("subscription", "chain", "grantee", "grantee_alias")
	require.NoError(t, err)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*MsgGrant)
	require.Equal(t, "grantee_alias", message.Grantee.Title)
	require.Equal(t, "granter_alias", message.Granter.Title)
}
