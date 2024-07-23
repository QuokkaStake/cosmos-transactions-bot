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

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMsgSendParse(t *testing.T) {
	t.Parallel()

	msg := &cosmosBankTypes.MsgSend{
		FromAddress: "from",
		ToAddress:   "to",
		Amount:      cosmosTypes.Coins{{Amount: cosmosTypes.NewInt(100), Denom: "ustake"}},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgSend(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed2, err2 := ParseMsgSend([]byte("aaa"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err2)
	require.Nil(t, parsed2)
}

func TestMsgSendBase(t *testing.T) {
	t.Parallel()

	msg := &cosmosBankTypes.MsgSend{
		FromAddress: "from",
		ToAddress:   "to",
		Amount:      cosmosTypes.Coins{{Amount: cosmosTypes.NewInt(100), Denom: "ustake"}},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgSend(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/cosmos.bank.v1beta1.MsgSend", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/cosmos.bank.v1beta1.MsgSend"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "from"),
		event.From(cosmosBankTypes.EventTypeTransfer, cosmosTypes.AttributeKeySender, "from"),
		event.From(cosmosBankTypes.EventTypeTransfer, cosmosBankTypes.AttributeKeyRecipient, "to"),
		event.From(cosmosBankTypes.EventTypeCoinSpent, cosmosBankTypes.AttributeKeySpender, "from"),
		event.From(cosmosBankTypes.EventTypeCoinReceived, cosmosBankTypes.AttributeKeyReceiver, "to"),
		event.From(cosmosBankTypes.EventTypeTransfer, cosmosTypes.AttributeKeyAmount, "100ustake"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "from"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgSendPopulate(t *testing.T) {
	t.Parallel()

	msg := &cosmosBankTypes.MsgSend{
		FromAddress: "from",
		ToAddress:   "to",
		Amount:      cosmosTypes.Coins{{Amount: cosmosTypes.NewInt(100000000), Denom: "uatom"}},
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
		AliasesPath: "path.toml",
	}

	parsed, err := ParseMsgSend(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	err = aliasManager.Set("subscription", "chain", "from", "from_alias")
	require.NoError(t, err)
	err = aliasManager.Set("subscription", "chain", "to", "to_alias")
	require.NoError(t, err)

	dataFetcher.Cache.Set("chain-id_price_uatom", 6.7)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	msgSend, _ := parsed.(*MsgSend)

	require.Len(t, msgSend.Amount, 1)
	require.Equal(t, "from_alias", msgSend.From.Title)
	require.Equal(t, "to_alias", msgSend.To.Title)

	first := msgSend.Amount[0]
	require.Equal(t, "100.00", fmt.Sprintf("%.2f", first.Value))
	require.Equal(t, "670.00", fmt.Sprintf("%.2f", first.PriceUSD))
	require.Equal(t, "atom", first.Denom.String())
}
