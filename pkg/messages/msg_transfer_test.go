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

	ibcTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMsgTransferParse(t *testing.T) {
	t.Parallel()

	msg := &ibcTypes.MsgTransfer{
		Token:    sdkTypes.Coin{Amount: sdkTypes.NewInt(100), Denom: "ustake"},
		Sender:   "sender",
		Receiver: "receiver",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgTransfer(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed2, err2 := ParseMsgTransfer([]byte("aaa"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err2)
	require.Nil(t, parsed2)
}

func TestMsgTransferBase(t *testing.T) {
	t.Parallel()

	msg := &ibcTypes.MsgTransfer{
		Token:    sdkTypes.Coin{Amount: sdkTypes.NewInt(100), Denom: "ustake"},
		Sender:   "sender",
		Receiver: "receiver",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgTransfer(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/ibc.applications.transfer.v1.MsgTransfer", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/ibc.applications.transfer.v1.MsgTransfer"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "sender"),
		event.From(ibcTypes.EventTypeTransfer, ibcTypes.AttributeKeyReceiver, "receiver"),
		event.From(ibcTypes.EventTypeTransfer, cosmosTypes.AttributeKeySender, "sender"),
		event.From(ibcTypes.EventTypeTransfer, cosmosTypes.AttributeKeyAmount, "100ustake"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgTransferPopulateLinks(t *testing.T) {
	t.Parallel()

	msg := &ibcTypes.MsgTransfer{
		Token:         sdkTypes.Coin{Amount: sdkTypes.NewInt(100000000), Denom: "uatom"},
		Sender:        "sender",
		Receiver:      "receiver",
		SourcePort:    "port",
		SourceChannel: "channel",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	config := &configPkg.AppConfig{
		Chains: configTypes.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
				Denoms: configTypes.DenomInfos{
					{Denom: "uatom", DisplayDenom: "atom", DenomExponent: 6, CoingeckoCurrency: "cosmos"},
				},
				Explorer: &configTypes.Explorer{WalletLinkPattern: "link %s"},
			},
			{
				Name:    "chain2",
				ChainID: "chain-id-2",
				Denoms: configTypes.DenomInfos{
					{Denom: "uatom", DisplayDenom: "atom", DenomExponent: 6, CoingeckoCurrency: "cosmos"},
				},
				Explorer: &configTypes.Explorer{WalletLinkPattern: "another link %s"},
			},
		},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	parsed, err := ParseMsgTransfer(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_channel_channel_port_port", "chain-id-2")

	err = aliasManager.Set("subscription", "chain", "sender", "sender_alias")
	require.NoError(t, err)
	err = aliasManager.Set("subscription", "chain2", "receiver", "receiver_alias")
	require.NoError(t, err)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*MsgTransfer)

	require.Equal(t, "sender_alias", message.Sender.Title)
	require.Equal(t, "receiver_alias", message.Receiver.Title)
	require.Equal(t, "link sender", message.Sender.Href)
	require.Equal(t, "another link receiver", message.Receiver.Href)
}

func TestMsgTransferPopulateNativeDenom(t *testing.T) {
	t.Parallel()

	msg := &ibcTypes.MsgTransfer{
		Token:         sdkTypes.Coin{Amount: sdkTypes.NewInt(100000000), Denom: "uatom"},
		Sender:        "sender",
		Receiver:      "receiver",
		SourcePort:    "port",
		SourceChannel: "channel",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

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
		AliasesPath: "path.yaml",
	}

	parsed, err := ParseMsgTransfer(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain-id_price_uatom", 6.7)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*MsgTransfer)
	require.Equal(t, "100.00", fmt.Sprintf("%.2f", message.Token.Value))
	require.Equal(t, "670.00", fmt.Sprintf("%.2f", message.Token.PriceUSD))
	require.Equal(t, "atom", message.Token.Denom.String())
}

func TestMsgTransferPopulateIbcDenomFetchDenomFailed(t *testing.T) {
	t.Parallel()

	msg := &ibcTypes.MsgTransfer{
		Token:         sdkTypes.Coin{Amount: sdkTypes.NewInt(100000000), Denom: "ibc/denom"},
		Sender:        "sender",
		Receiver:      "receiver",
		SourcePort:    "port",
		SourceChannel: "channel",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	config := &configPkg.AppConfig{
		Chains: configTypes.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
			},
			{
				Name:    "chain2",
				ChainID: "remote-chain",
				Denoms: configTypes.DenomInfos{
					{Denom: "uatom", DisplayDenom: "atom", DenomExponent: 6, CoingeckoCurrency: "cosmos"},
				},
			},
		},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	parsed, err := ParseMsgTransfer(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_denom_trace_denom", nil)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*MsgTransfer)
	require.Equal(t, "100000000.00", fmt.Sprintf("%.2f", message.Token.Value))
	require.Nil(t, message.Token.PriceUSD)
	require.Equal(t, "ibc/denom", message.Token.Denom.String())
}

func TestMsgTransferPopulateIbcDenomFetchRemoteChainFailed(t *testing.T) {
	t.Parallel()

	msg := &ibcTypes.MsgTransfer{
		Token:         sdkTypes.Coin{Amount: sdkTypes.NewInt(100000000), Denom: "ibc/denom"},
		Sender:        "sender",
		Receiver:      "receiver",
		SourcePort:    "port",
		SourceChannel: "channel",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	config := &configPkg.AppConfig{
		Chains: configTypes.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
			},
			{
				Name:    "chain2",
				ChainID: "remote-chain",
				Denoms: configTypes.DenomInfos{
					{Denom: "uatom", DisplayDenom: "atom", DenomExponent: 6, CoingeckoCurrency: "cosmos"},
				},
			},
		},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	parsed, err := ParseMsgTransfer(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_denom_trace_denom", &ibcTypes.DenomTrace{
		Path:      "path",
		BaseDenom: "uatom",
	})
	dataFetcher.Cache.Set("chain_channel_channel_port_port", nil)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*MsgTransfer)
	require.Equal(t, "100000000.00", fmt.Sprintf("%.2f", message.Token.Value))
	require.Nil(t, message.Token.PriceUSD)
	require.Equal(t, "uatom", message.Token.Denom.String())
}

func TestMsgTransferPopulateIbcDenom(t *testing.T) {
	t.Parallel()

	msg := &ibcTypes.MsgTransfer{
		Token:         sdkTypes.Coin{Amount: sdkTypes.NewInt(100000000), Denom: "ibc/denom"},
		Sender:        "sender",
		Receiver:      "receiver",
		SourcePort:    "port",
		SourceChannel: "channel",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	config := &configPkg.AppConfig{
		Chains: configTypes.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
			},
			{
				Name:    "chain2",
				ChainID: "remote-chain",
				Denoms: configTypes.DenomInfos{
					{Denom: "uatom", DisplayDenom: "atom", DenomExponent: 6, CoingeckoCurrency: "cosmos"},
				},
			},
		},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	parsed, err := ParseMsgTransfer(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_denom_trace_denom", &ibcTypes.DenomTrace{
		Path:      "path",
		BaseDenom: "uatom",
	})
	dataFetcher.Cache.Set("chain_channel_channel_port_port", "remote-chain")
	dataFetcher.Cache.Set("remote-chain_price_uatom", 6.7)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*MsgTransfer)
	require.Equal(t, "100.00", fmt.Sprintf("%.2f", message.Token.Value))
	require.Equal(t, "670.00", fmt.Sprintf("%.2f", message.Token.PriceUSD))
	require.Equal(t, "atom", message.Token.Denom.String())
}
