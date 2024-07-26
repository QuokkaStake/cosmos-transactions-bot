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

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	ibcTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcChannelTypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMsgAcknowledgementParseFail(t *testing.T) {
	t.Parallel()

	parsed, err := ParseMsgAcknowledgement([]byte("bytes"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err)
	require.Nil(t, parsed)
}

func TestMsgAcknowledgementParsePacketFail(t *testing.T) {
	t.Parallel()

	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "dst_channel",
		DestinationPort:    "dst_port",
		Data:               []byte("invalid"),
	}

	msg := &ibcChannelTypes.MsgAcknowledgement{
		Signer: "signer",
		Packet: packet,
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgAcknowledgement(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err)
	require.Nil(t, parsed)
}

func TestMsgAcknowledgementParseOk(t *testing.T) {
	t.Parallel()

	msgInternal := &ibcTypes.FungibleTokenPacketData{
		Denom:    "uatom",
		Amount:   "100",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     "memo",
	}

	msgInternalBytes, err := ibcTypes.ModuleCdc.MarshalJSON(msgInternal)
	require.NoError(t, err)

	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "dst_channel",
		DestinationPort:    "dst_port",
		Data:               msgInternalBytes,
	}

	msg := &ibcChannelTypes.MsgAcknowledgement{
		Signer: "signer",
		Packet: packet,
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgAcknowledgement(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)
}

func TestMsgAcknowledgementBase(t *testing.T) {
	t.Parallel()

	msgInternal := &ibcTypes.FungibleTokenPacketData{
		Denom:    "uatom",
		Amount:   "100",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     "memo",
	}

	msgInternalBytes, err := ibcTypes.ModuleCdc.MarshalJSON(msgInternal)
	require.NoError(t, err)

	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "dst_channel",
		DestinationPort:    "dst_port",
		Data:               msgInternalBytes,
	}

	msg := &ibcChannelTypes.MsgAcknowledgement{
		Signer: "signer",
		Packet: packet,
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgAcknowledgement(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/ibc.core.channel.v1.MsgAcknowledgement", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/ibc.core.channel.v1.MsgAcknowledgement"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "signer"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgAcknowledgementPopulate(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: configTypes.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
				Denoms: configTypes.DenomInfos{
					{Denom: "uatom", DisplayDenom: "atom", DenomExponent: 6, CoingeckoCurrency: "cosmos"},
				},
			},
		},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.toml",
	}

	msgInternal := &ibcTypes.FungibleTokenPacketData{
		Denom:    "uatom",
		Amount:   "100",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     "memo",
	}

	msgInternalBytes, err := ibcTypes.ModuleCdc.MarshalJSON(msgInternal)
	require.NoError(t, err)

	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "dst_channel",
		DestinationPort:    "dst_port",
		Data:               msgInternalBytes,
	}

	msg := &ibcChannelTypes.MsgAcknowledgement{
		Signer: "signer",
		Packet: packet,
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgAcknowledgement(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	err = aliasManager.Set("subscription", "chain", "sender", "sender_alias")
	require.NoError(t, err)
	err = aliasManager.Set("subscription", "chain", "signer", "signer_alias")
	require.NoError(t, err)
	err = aliasManager.Set("subscription", "chain", "receiver", "receiver_alias")
	require.NoError(t, err)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, ok := parsed.(*MsgAcknowledgement)
	require.True(t, ok)
	require.Equal(t, "sender_alias", message.Sender.Title)
	require.Equal(t, "signer_alias", message.Signer.Title)
}
