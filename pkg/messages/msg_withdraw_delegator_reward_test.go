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

func TestMsgWithdrawDelegatorRewardParse(t *testing.T) {
	t.Parallel()

	msg := &cosmosDistributionTypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: "delegator",
		ValidatorAddress: "validator",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgWithdrawDelegatorReward(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed2, err2 := ParseMsgWithdrawDelegatorReward([]byte("aaa"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err2)
	require.Nil(t, parsed2)
}

func TestMsgWithdrawDelegatorRewardBase(t *testing.T) {
	t.Parallel()

	msg := &cosmosDistributionTypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: "delegator",
		ValidatorAddress: "validator",
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgWithdrawDelegatorReward(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "delegator"),
		event.From(cosmosDistributionTypes.EventTypeWithdrawRewards, cosmosDistributionTypes.AttributeKeyValidator, "validator"),
		event.From(cosmosDistributionTypes.EventTypeWithdrawRewards, cosmosDistributionTypes.AttributeKeyDelegator, "delegator"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgWithdrawDelegatorRewardPopulate(t *testing.T) {
	t.Parallel()

	msg := &cosmosDistributionTypes.MsgWithdrawDelegatorReward{
		DelegatorAddress: "delegator",
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
		AliasesPath: "path.yaml",
	}

	parsed, err := ParseMsgWithdrawDelegatorReward(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	err = aliasManager.Set("subscription", "chain", "delegator", "delegator_alias")
	require.NoError(t, err)

	dataFetcher.Cache.Set("chain_rewards_delegator_validator_100", []responses.Reward{
		{Amount: "100000000", Denom: "uatom"},
	})
	dataFetcher.Cache.Set("chain_validator_validator", &responses.Validator{
		OperatorAddress: "test",
		Description:     responses.ValidatorDescription{Moniker: "Validator Moniker"},
	})
	dataFetcher.Cache.Set("chain-id_price_uatom", 6.7)

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*MsgWithdrawDelegatorReward)

	require.Equal(t, "delegator_alias", message.DelegatorAddress.Title)
	require.Equal(t, "Validator Moniker", message.ValidatorAddress.Title)

	require.Len(t, message.Amount, 1)
	require.Equal(t, "100.00", fmt.Sprintf("%.2f", message.Amount[0].Value))
	require.Equal(t, "670.00", fmt.Sprintf("%.2f", message.Amount[0].PriceUSD))
	require.Equal(t, "atom", message.Amount[0].Denom.String())
}
