package messages

import (
	aliasManagerPkg "main/pkg/alias_manager"
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
	"main/pkg/data_fetcher"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"
	packet2 "main/pkg/messages/packet"
	"main/pkg/metrics"
	"main/pkg/types"
	"main/pkg/types/event"
	"testing"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcChannelTypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMsgRecvPacketParseFail(t *testing.T) {
	t.Parallel()

	parsed, err := ParseMsgRecvPacket([]byte("bytes"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err)
	require.Nil(t, parsed)
}

func TestMsgRecvPacketParsePacketFail(t *testing.T) {
	t.Parallel()

	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "dst_channel",
		DestinationPort:    "dst_port",
		Data:               []byte("invalid"),
	}

	msg := &ibcChannelTypes.MsgRecvPacket{
		Signer: "signer",
		Packet: packet,
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgRecvPacket(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err)
	require.Nil(t, parsed)
}

func TestMsgRecvPacketParseOk(t *testing.T) {
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

	msg := &ibcChannelTypes.MsgRecvPacket{
		Signer: "signer",
		Packet: packet,
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgRecvPacket(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)
}

func TestMsgRecvPacketBase(t *testing.T) {
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

	msg := &ibcChannelTypes.MsgRecvPacket{
		Signer: "signer",
		Packet: packet,
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgRecvPacket(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/ibc.core.channel.v1.MsgRecvPacket", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/ibc.core.channel.v1.MsgRecvPacket"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "signer"),
		event.From(ibcTypes.EventTypePacket, cosmosBankTypes.AttributeKeyReceiver, "receiver"),
		event.From(ibcTypes.EventTypePacket, cosmosTypes.AttributeKeyAmount, "100"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "sender"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgRecvPacketPopulate(t *testing.T) {
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

	msg := &ibcChannelTypes.MsgRecvPacket{
		Signer: "signer",
		Packet: packet,
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgRecvPacket(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	err = aliasManager.Set("subscription", "chain", "signer", "signer_alias")
	require.NoError(t, err)
	err = aliasManager.Set("subscription", "chain", "receiver", "receiver_alias")
	require.NoError(t, err)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, ok := parsed.(*MsgRecvPacket)
	require.True(t, ok)
	require.Equal(t, "signer_alias", message.Signer.Title)

	packetParsed, ok := message.Packet.(*packet2.FungibleTokenPacket)
	require.True(t, ok)
	require.Equal(t, "receiver_alias", packetParsed.Receiver.Title)
}
