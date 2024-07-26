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

	types2 "github.com/cosmos/cosmos-sdk/codec/types"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	cosmosAuthzTypes "github.com/cosmos/cosmos-sdk/x/authz"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMsgExecParse(t *testing.T) {
	t.Parallel()

	msgInternal := &cosmosDistributionTypes.MsgSetWithdrawAddress{WithdrawAddress: "address"}
	msgInternalBytes, err := proto.Marshal(msgInternal)
	require.NoError(t, err)

	msg := &cosmosAuthzTypes.MsgExec{
		Grantee: "grantee",
		Msgs: []*types2.Any{
			{TypeUrl: "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress", Value: msgInternalBytes},
		},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgExec(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed2, err2 := ParseMsgExec([]byte("aaa"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err2)
	require.Nil(t, parsed2)
}

func TestMsgExecBase(t *testing.T) {
	t.Parallel()

	msgInternal := &cosmosDistributionTypes.MsgSetWithdrawAddress{DelegatorAddress: "delegator"}
	msgInternalBytes, err := proto.Marshal(msgInternal)
	require.NoError(t, err)

	msgSetAddr, err := ParseMsgSetWithdrawAddress(msgInternalBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)

	msg := &cosmosAuthzTypes.MsgExec{
		Grantee: "grantee",
		Msgs: []*types2.Any{
			{TypeUrl: "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress", Value: msgInternalBytes},
		},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgExec(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/cosmos.authz.v1beta1.MsgExec", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/cosmos.authz.v1beta1.MsgExec"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "grantee"),
	}, values)

	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	parsed.AddParsedMessage(msgSetAddr)
	require.Len(t, parsed.GetParsedMessages(), 1)
	require.Len(t, parsed.GetRawMessages(), 1)
}

func TestMsgExecPopulate(t *testing.T) {
	t.Parallel()

	msgInternal := &cosmosDistributionTypes.MsgSetWithdrawAddress{
		DelegatorAddress: "delegator",
		WithdrawAddress:  "withdraw",
	}
	msgInternalBytes, err := proto.Marshal(msgInternal)
	require.NoError(t, err)

	msgSetAddr, err := ParseMsgSetWithdrawAddress(msgInternalBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)

	msg := &cosmosAuthzTypes.MsgExec{
		Grantee: "grantee",
		Msgs: []*types2.Any{
			{TypeUrl: "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress", Value: msgInternalBytes},
		},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	config := &configPkg.AppConfig{
		Chains:      configTypes.Chains{{Name: "chain"}},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.toml",
	}

	parsed, err := ParseMsgExec(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed.AddParsedMessage(msgSetAddr)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	err = aliasManager.Set("subscription", "chain", "delegator", "delegator_alias")
	require.NoError(t, err)
	err = aliasManager.Set("subscription", "chain", "grantee", "grantee_alias")
	require.NoError(t, err)
	err = aliasManager.Set("subscription", "chain", "withdraw", "withdraw_alias")
	require.NoError(t, err)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*MsgExec)

	require.Equal(t, "grantee_alias", message.Grantee.Title)
	require.Len(t, parsed.GetParsedMessages(), 1)
	internal, ok := parsed.GetParsedMessages()[0].(*MsgSetWithdrawAddress)
	require.True(t, ok)
	require.Equal(t, "delegator_alias", internal.DelegatorAddress.Title)
	require.Equal(t, "withdraw_alias", internal.WithdrawAddress.Title)
}

func TestMsgExecLabel(t *testing.T) {
	t.Parallel()

	msgInternal := &cosmosDistributionTypes.MsgSetWithdrawAddress{DelegatorAddress: "delegator"}
	msgInternalBytes, err := proto.Marshal(msgInternal)
	require.NoError(t, err)

	msg := &cosmosAuthzTypes.MsgExec{
		Grantee: "grantee",
		Msgs: []*types2.Any{
			{TypeUrl: "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress", Value: msgInternalBytes},
		},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgExec(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	message, _ := parsed.(*MsgExec)

	require.Equal(t, "1, 1 skipped", message.GetMessagesLabel())

	msgSetAddr, err := ParseMsgSetWithdrawAddress(msgInternalBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)

	parsed.AddParsedMessage(msgSetAddr)
	require.Equal(t, "1", message.GetMessagesLabel())
}
