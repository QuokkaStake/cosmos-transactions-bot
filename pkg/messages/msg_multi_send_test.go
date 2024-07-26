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

func TestMsgMultiSendParse(t *testing.T) {
	t.Parallel()

	msg := &cosmosBankTypes.MsgMultiSend{
		Inputs: []cosmosBankTypes.Input{
			{
				Address: "from",
				Coins:   cosmosTypes.Coins{{Amount: cosmosTypes.NewInt(100), Denom: "ustake"}},
			},
		},
		Outputs: []cosmosBankTypes.Output{
			{
				Address: "from",
				Coins:   cosmosTypes.Coins{{Amount: cosmosTypes.NewInt(100), Denom: "ustake"}},
			},
		},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgMultiSend(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed2, err2 := ParseMsgMultiSend([]byte("aaa"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err2)
	require.Nil(t, parsed2)
}

func TestMsgMultiSendBase(t *testing.T) {
	t.Parallel()

	msg := &cosmosBankTypes.MsgMultiSend{
		Inputs: []cosmosBankTypes.Input{
			{
				Address: "from",
				Coins:   cosmosTypes.Coins{{Amount: cosmosTypes.NewInt(100), Denom: "ustake"}},
			},
		},
		Outputs: []cosmosBankTypes.Output{
			{
				Address: "from",
				Coins:   cosmosTypes.Coins{{Amount: cosmosTypes.NewInt(100), Denom: "ustake"}},
			},
		},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgMultiSend(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/cosmos.bank.v1beta1.MsgMultiSend", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/cosmos.bank.v1beta1.MsgMultiSend"),
		event.From(cosmosBankTypes.EventTypeTransfer, cosmosBankTypes.AttributeKeySpender, "from"),
		event.From(cosmosBankTypes.EventTypeCoinSpent, cosmosBankTypes.AttributeKeySpender, "from"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "from"),
		event.From(cosmosBankTypes.EventTypeTransfer, cosmosTypes.AttributeKeyAmount, "100ustake"),
		event.From(cosmosBankTypes.EventTypeTransfer, cosmosBankTypes.AttributeKeyRecipient, "from"),
		event.From(cosmosBankTypes.EventTypeCoinReceived, cosmosBankTypes.AttributeKeyReceiver, "from"),
		event.From(cosmosBankTypes.EventTypeTransfer, cosmosTypes.AttributeKeyAmount, "100ustake"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgMultiSendPopulate(t *testing.T) {
	t.Parallel()

	msg := &cosmosBankTypes.MsgMultiSend{
		Inputs: []cosmosBankTypes.Input{
			{
				Address: "from",
				Coins:   cosmosTypes.Coins{{Amount: cosmosTypes.NewInt(100000000), Denom: "uatom"}},
			},
		},
		Outputs: []cosmosBankTypes.Output{
			{
				Address: "to",
				Coins:   cosmosTypes.Coins{{Amount: cosmosTypes.NewInt(100000000), Denom: "uatom"}},
			},
		},
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

	parsed, err := ParseMsgMultiSend(msgBytes, config.Chains[0], 100)
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

	msgSend, _ := parsed.(*MsgMultiSend)

	require.Len(t, msgSend.Inputs, 1)
	require.Len(t, msgSend.Outputs, 1)
	require.Equal(t, "from_alias", msgSend.Inputs[0].Address.Title)
	require.Equal(t, "to_alias", msgSend.Outputs[0].Address.Title)

	require.Len(t, msgSend.Inputs[0].Amount, 1)
	require.Len(t, msgSend.Outputs[0].Amount, 1)

	firstInput := msgSend.Inputs[0].Amount[0]
	require.Equal(t, "100.00", fmt.Sprintf("%.2f", firstInput.Value))
	require.Equal(t, "670.00", fmt.Sprintf("%.2f", firstInput.PriceUSD))
	require.Equal(t, "atom", firstInput.Denom.String())

	firstOutput := msgSend.Outputs[0].Amount[0]
	require.Equal(t, "100.00", fmt.Sprintf("%.2f", firstOutput.Value))
	require.Equal(t, "670.00", fmt.Sprintf("%.2f", firstOutput.PriceUSD))
	require.Equal(t, "atom", firstOutput.Denom.String())
}
