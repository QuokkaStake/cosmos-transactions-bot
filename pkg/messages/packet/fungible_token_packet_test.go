package packet

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
	"main/pkg/types/responses"
	"testing"

	ibcTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcChannelTypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

func TestFungibleTokenPacketParse(t *testing.T) {
	t.Parallel()

	msg := ibcTypes.FungibleTokenPacketData{
		Denom:    "uatom",
		Amount:   "100",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     "memo",
	}
	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "dst_channel",
		DestinationPort:    "dst_port",
	}

	parsed := ParseFungibleTokenPacket(msg, packet, &configTypes.Chain{Name: "chain"})
	require.NotNil(t, parsed)
}

func TestFungibleTokenPacketBase(t *testing.T) {
	t.Parallel()

	msg := ibcTypes.FungibleTokenPacketData{
		Denom:    "uatom",
		Amount:   "100",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     "memo",
	}
	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "dst_channel",
		DestinationPort:    "dst_port",
	}

	parsed := ParseFungibleTokenPacket(msg, packet, &configTypes.Chain{Name: "chain"})
	require.NotNil(t, parsed)

	require.Equal(t, "FungibleTokenPacket", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(ibcTypes.EventTypePacket, cosmosBankTypes.AttributeKeyReceiver, "receiver"),
		event.From(ibcTypes.EventTypePacket, cosmosTypes.AttributeKeyAmount, "100"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "sender"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestFungibleTokenPacketPopulateLinks(t *testing.T) {
	t.Parallel()

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

	msg := ibcTypes.FungibleTokenPacketData{
		Denom:    "uatom",
		Amount:   "100",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     "memo",
	}
	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "channel",
		DestinationPort:    "port",
	}

	parsed := ParseFungibleTokenPacket(msg, packet, config.Chains[0])
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_channel_channel_port_port", "chain-id-2")

	err := aliasManager.Set("subscription", "chain2", "sender", "sender_alias")
	require.NoError(t, err)
	err = aliasManager.Set("subscription", "chain", "receiver", "receiver_alias")
	require.NoError(t, err)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*FungibleTokenPacket)

	require.Equal(t, "sender_alias", message.Sender.Title)
	require.Equal(t, "receiver_alias", message.Receiver.Title)
	require.Equal(t, "another link sender", message.Sender.Href)
	require.Equal(t, "link receiver", message.Receiver.Href)
}

func TestMsgTransferPopulateNativeDenomNotFetched(t *testing.T) {
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
		AliasesPath: "path.yaml",
	}

	msg := ibcTypes.FungibleTokenPacketData{
		Denom:    "uatom",
		Amount:   "100",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     "memo",
	}
	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "channel",
		DestinationPort:    "port",
	}

	parsed := ParseFungibleTokenPacket(msg, packet, config.Chains[0])
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_channel_channel_port_port", nil)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*FungibleTokenPacket)
	require.Equal(t, "100.00", fmt.Sprintf("%.2f", message.Token.Value))
	require.Nil(t, message.Token.PriceUSD)
	require.Equal(t, "uatom", message.Token.Denom.String())
}

func TestMsgTransferPopulateNotNativeDenom(t *testing.T) {
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
		AliasesPath: "path.yaml",
	}

	msg := ibcTypes.FungibleTokenPacketData{
		Denom:    "transfer/channel-0/uatom",
		Amount:   "100000000",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     "memo",
	}
	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "channel",
		DestinationPort:    "port",
	}

	parsed := ParseFungibleTokenPacket(msg, packet, config.Chains[0])
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain-id_price_uatom", 6.7)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*FungibleTokenPacket)
	require.Equal(t, "100.00", fmt.Sprintf("%.2f", message.Token.Value))
	require.Equal(t, "670.00", fmt.Sprintf("%.2f", message.Token.PriceUSD))
	require.Equal(t, "atom", message.Token.Denom.String())
}

func TestMsgTransferPopulateNativeDenomNotFoundLocally(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: configTypes.Chains{
			{
				Name:    "chain",
				ChainID: "chain-id",
			},
		},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.yaml",
	}

	msg := ibcTypes.FungibleTokenPacketData{
		Denom:    "uatom",
		Amount:   "100000000",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     "memo",
	}
	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "channel",
		DestinationPort:    "port",
	}

	parsed := ParseFungibleTokenPacket(msg, packet, config.Chains[0])
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_denom_trace_denom", &ibcTypes.DenomTrace{
		Path:      "port/channel",
		BaseDenom: "uatom",
	})
	dataFetcher.Cache.Set("cosmos_directory_chains", responses.CosmosDirectoryChains{
		{
			ChainID: "remote-chain",
			Assets: []responses.CosmosDirectoryAsset{
				{
					Denom:       "uatom",
					CoingeckoID: "cosmos",
					Base:        responses.CosmosDirectoryAssetDenomInfo{Denom: "uatom"},
					Display:     responses.CosmosDirectoryAssetDenomInfo{Denom: "atom", Exponent: 6},
				},
			},
		},
	})
	dataFetcher.Cache.Set("remote-chain_price_uatom", 6.7)
	dataFetcher.Cache.Set("chain_channel_channel_port_port", "remote-chain")

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*FungibleTokenPacket)
	require.Equal(t, "100.00", fmt.Sprintf("%.2f", message.Token.Value))
	require.Equal(t, "670.00", fmt.Sprintf("%.2f", message.Token.PriceUSD))
	require.Equal(t, "atom", message.Token.Denom.String())
}
