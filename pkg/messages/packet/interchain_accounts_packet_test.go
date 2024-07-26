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
	"main/pkg/types/amount"
	"math/big"
	"testing"

	icaTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"

	types2 "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestInterchainAccountPacketParseFail(t *testing.T) {
	t.Parallel()

	packet := icaTypes.InterchainAccountPacketData{
		Type: icaTypes.EXECUTE_TX,
		Data: []byte("random"),
	}

	parsed, err := ParseInterchainAccountsPacket(packet, &configTypes.Chain{Name: "chain"})
	require.Error(t, err)
	require.Nil(t, parsed)
}

func TestInterchainAccountPacketParseOk(t *testing.T) {
	t.Parallel()

	txInternal := &icaTypes.CosmosTx{
		Messages: []*types2.Any{},
	}
	txBytes, err := proto.Marshal(txInternal)
	require.NoError(t, err)

	packet := icaTypes.InterchainAccountPacketData{
		Type: icaTypes.EXECUTE_TX,
		Data: txBytes,
	}

	parsed, err := ParseInterchainAccountsPacket(packet, &configTypes.Chain{Name: "chain"})
	require.NoError(t, err)
	require.NotNil(t, parsed)
}

func TestInterchainAccountPacketBase(t *testing.T) {
	t.Parallel()

	txInternal := &icaTypes.CosmosTx{
		Messages: []*types2.Any{},
	}
	txBytes, err := proto.Marshal(txInternal)
	require.NoError(t, err)

	packet := icaTypes.InterchainAccountPacketData{
		Type: icaTypes.EXECUTE_TX,
		Data: txBytes,
	}

	parsed, err := ParseInterchainAccountsPacket(packet, &configTypes.Chain{Name: "chain"})
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "InterchainAccountsPacket", parsed.Type())
	values := parsed.GetValues()
	require.Empty(t, values)

	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
	parsed.AddParsedMessage(&FungibleTokenPacket{
		Receiver: &configTypes.Link{Value: "receiver"},
		Sender:   &configTypes.Link{Value: "sender"},
		Token:    &amount.Amount{Value: big.NewFloat(1), Denom: "ustake"},
	})

	values2 := parsed.GetValues()
	require.NotEmpty(t, values2)
}

func TestInterchainAccountPacketLabel(t *testing.T) {
	t.Parallel()

	txInternal := &icaTypes.CosmosTx{
		Messages: []*types2.Any{
			{TypeUrl: "/cosmos.bank.v1beta1.MsgSend", Value: nil},
		},
	}
	txBytes, err := proto.Marshal(txInternal)
	require.NoError(t, err)

	packet := icaTypes.InterchainAccountPacketData{
		Type: icaTypes.EXECUTE_TX,
		Data: txBytes,
	}

	parsed, err := ParseInterchainAccountsPacket(packet, &configTypes.Chain{Name: "chain"})
	require.NoError(t, err)
	require.NotNil(t, parsed)

	message, ok := parsed.(*InterchainAccountsPacket)
	require.True(t, ok)

	require.Equal(t, "1, 1 skipped", message.GetMessagesLabel())

	parsed.AddParsedMessage(nil)
	require.Equal(t, "1", message.GetMessagesLabel())
}

func TestInterchainAccountPacketPopulate(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: configTypes.Chains{{
			Name:    "chain",
			ChainID: "chain-id",
			Denoms: configTypes.DenomInfos{
				{Denom: "uatom", DisplayDenom: "atom", DenomExponent: 6, CoingeckoCurrency: "cosmos"},
			},
		}},
		Metrics:     configPkg.MetricsConfig{Enabled: false},
		AliasesPath: "path.toml",
	}

	txInternal := &icaTypes.CosmosTx{
		Messages: []*types2.Any{},
	}
	txBytes, err := proto.Marshal(txInternal)
	require.NoError(t, err)

	packet := icaTypes.InterchainAccountPacketData{
		Type: icaTypes.EXECUTE_TX,
		Data: txBytes,
	}

	parsed, err := ParseInterchainAccountsPacket(packet, &configTypes.Chain{Name: "chain"})
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed.AddParsedMessage(&FungibleTokenPacket{
		Receiver: &configTypes.Link{Value: "receiver"},
		Sender:   &configTypes.Link{Value: "sender"},
		Token:    &amount.Amount{Value: big.NewFloat(100000000), Denom: "transfer/channel-0/uatom"},
		Chain:    config.Chains[0],
	})

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_channel_channel_port_port", "chain-id")
	dataFetcher.Cache.Set("chain-id_price_uatom", 6.7)
	parsed.GetAdditionalData(dataFetcher, "subscription")

	parsedMessages := parsed.GetParsedMessages()
	require.Len(t, parsedMessages, 1)
	messageInternal, ok := parsedMessages[0].(*FungibleTokenPacket)
	require.True(t, ok)

	require.Equal(t, "100.00", fmt.Sprintf("%.2f", messageInternal.Token.Value))
	require.Equal(t, "670.00", fmt.Sprintf("%.2f", messageInternal.Token.PriceUSD))
	require.Equal(t, "atom", messageInternal.Token.Denom.String())
}
