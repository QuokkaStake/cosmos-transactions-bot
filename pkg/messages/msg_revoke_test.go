package messages

import (
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

	cosmosAuthzTypes "github.com/cosmos/cosmos-sdk/x/authz"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMsgRevokeParse(t *testing.T) {
	t.Parallel()

	msg := &cosmosAuthzTypes.MsgRevoke{
		Granter: "granter",
		Grantee: "grantee",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgRevoke(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed2, err2 := ParseMsgRevoke([]byte("aaa"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err2)
	require.Nil(t, parsed2)
}

func TestMsgRevokeBase(t *testing.T) {
	t.Parallel()

	msg := &cosmosAuthzTypes.MsgRevoke{
		Granter: "granter",
		Grantee: "grantee",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgRevoke(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/cosmos.authz.v1beta1.MsgRevoke", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/cosmos.authz.v1beta1.MsgRevoke"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "grantee"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgRevokePopulate(t *testing.T) {
	t.Parallel()

	msg := &cosmosAuthzTypes.MsgRevoke{
		Granter: "granter",
		Grantee: "grantee",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	config := &configPkg.AppConfig{
		Chains:      configTypes.Chains{{Name: "chain"}},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	parsed, err := ParseMsgRevoke(msgBytes, config.Chains[0], 100)
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

	message, _ := parsed.(*MsgRevoke)

	require.Equal(t, "granter_alias", message.Granter.Title)
	require.Equal(t, "grantee_alias", message.Grantee.Title)
}
