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

	ibcClientTypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateClientParse(t *testing.T) {
	t.Parallel()

	msg := &ibcClientTypes.MsgUpdateClient{
		ClientId: "client",
		Signer:   "signer",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgUpdateClient(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed2, err2 := ParseMsgUpdateClient([]byte("aaa"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err2)
	require.Nil(t, parsed2)
}

func TestMsgUpdateClientBase(t *testing.T) {
	t.Parallel()

	msg := &ibcClientTypes.MsgUpdateClient{
		ClientId: "client",
		Signer:   "signer",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgUpdateClient(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/ibc.core.client.v1.MsgUpdateClient", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/ibc.core.client.v1.MsgUpdateClient"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "signer"),
		event.From(ibcClientTypes.EventTypeUpdateClient, ibcClientTypes.AttributeKeyClientID, "client"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgUpdateClientPopulate(t *testing.T) {
	t.Parallel()

	msg := &ibcClientTypes.MsgUpdateClient{
		ClientId: "client",
		Signer:   "signer",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	config := &configPkg.AppConfig{
		Chains:      configTypes.Chains{{Name: "chain"}},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	parsed, err := ParseMsgUpdateClient(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	err = aliasManager.Set("subscription", "chain", "signer", "signer_alias")
	require.NoError(t, err)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*MsgUpdateClient)

	require.Equal(t, "signer_alias", message.Signer.Title)
}
