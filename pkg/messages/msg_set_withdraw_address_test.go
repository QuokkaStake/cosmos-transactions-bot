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

	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMsgSetWithdrawAddressParse(t *testing.T) {
	t.Parallel()

	msg := &cosmosDistributionTypes.MsgSetWithdrawAddress{
		DelegatorAddress: "delegator",
		WithdrawAddress:  "withdraw",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgSetWithdrawAddress(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed2, err2 := ParseMsgSetWithdrawAddress([]byte("aaa"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err2)
	require.Nil(t, parsed2)
}

func TestMsgSetWithdrawAddressBase(t *testing.T) {
	t.Parallel()

	msg := &cosmosDistributionTypes.MsgSetWithdrawAddress{
		DelegatorAddress: "delegator",
		WithdrawAddress:  "withdraw",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgSetWithdrawAddress(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "delegator"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgSetWithdrawAddressPopulate(t *testing.T) {
	t.Parallel()

	msg := &cosmosDistributionTypes.MsgSetWithdrawAddress{
		DelegatorAddress: "delegator",
		WithdrawAddress:  "withdraw",
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

	parsed, err := ParseMsgSetWithdrawAddress(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	err = aliasManager.Set("subscription", "chain", "delegator", "delegator_alias")
	require.NoError(t, err)
	err = aliasManager.Set("subscription", "chain", "withdraw", "withdraw_alias")
	require.NoError(t, err)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	msgSend, _ := parsed.(*MsgSetWithdrawAddress)

	require.Equal(t, "delegator_alias", msgSend.DelegatorAddress.Title)
	require.Equal(t, "withdraw_alias", msgSend.WithdrawAddress.Title)
}
