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
	"main/pkg/types/responses"
	"testing"

	cosmosDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMsgWithdrawValidatorCommissionParse(t *testing.T) {
	t.Parallel()

	msg := &cosmosDistributionTypes.MsgWithdrawValidatorCommission{
		ValidatorAddress: "validator",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgWithdrawValidatorCommission(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed2, err2 := ParseMsgWithdrawValidatorCommission([]byte("aaa"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err2)
	require.Nil(t, parsed2)
}

func TestMsgWithdrawValidatorCommissionBase(t *testing.T) {
	t.Parallel()

	msg := &cosmosDistributionTypes.MsgWithdrawValidatorCommission{
		ValidatorAddress: "validator",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgWithdrawValidatorCommission(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"),
		event.From(cosmosDistributionTypes.EventTypeWithdrawCommission, cosmosDistributionTypes.AttributeKeyValidator, "validator"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgWithdrawValidatorCommissionPopulate(t *testing.T) {
	t.Parallel()

	msg := &cosmosDistributionTypes.MsgWithdrawValidatorCommission{
		ValidatorAddress: "validator",
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

	parsed, err := ParseMsgWithdrawValidatorCommission(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	dataFetcher.Cache.Set("chain_commission_validator_100", []responses.Commission{
		{Amount: "100000000", Denom: "uatom"},
	})
	dataFetcher.Cache.Set("chain_validator_validator", &responses.Validator{
		OperatorAddress: "test",
		Description:     responses.ValidatorDescription{Moniker: "Validator Moniker"},
	})
	dataFetcher.Cache.Set("chain-id_price_uatom", 6.7)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*MsgWithdrawValidatorCommission)

	require.Equal(t, "Validator Moniker", message.ValidatorAddress.Title)
	require.Len(t, message.Amount, 1)
	require.Equal(t, "100.00", fmt.Sprintf("%.2f", message.Amount[0].Value))
	require.Equal(t, "670.00", fmt.Sprintf("%.2f", message.Amount[0].PriceUSD))
	require.Equal(t, "atom", message.Amount[0].Denom.String())
}
