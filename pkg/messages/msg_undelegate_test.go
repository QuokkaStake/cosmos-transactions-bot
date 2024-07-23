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
	"time"

	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMsgUndelegateParse(t *testing.T) {
	t.Parallel()

	msg := &cosmosStakingTypes.MsgUndelegate{
		DelegatorAddress: "delegator",
		ValidatorAddress: "validator",
		Amount:           cosmosTypes.Coin{Amount: cosmosTypes.NewInt(100), Denom: "ustake"},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgUndelegate(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed2, err2 := ParseMsgUndelegate([]byte("aaa"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err2)
	require.Nil(t, parsed2)
}

func TestMsgUndelegateBase(t *testing.T) {
	t.Parallel()

	msg := &cosmosStakingTypes.MsgUndelegate{
		DelegatorAddress: "delegator",
		ValidatorAddress: "validator",
		Amount:           cosmosTypes.Coin{Amount: cosmosTypes.NewInt(100), Denom: "ustake"},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgUndelegate(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/cosmos.staking.v1beta1.MsgUndelegate", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/cosmos.staking.v1beta1.MsgUndelegate"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "delegator"),
		event.From(cosmosStakingTypes.EventTypeUnbond, cosmosStakingTypes.AttributeKeyValidator, "validator"),
		event.From(cosmosStakingTypes.EventTypeUnbond, cosmosStakingTypes.AttributeKeyDelegator, "delegator"),
		event.From(cosmosStakingTypes.EventTypeUnbond, cosmosTypes.AttributeKeyAmount, "100ustake"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgUndelegatePopulate(t *testing.T) {
	t.Parallel()

	msg := &cosmosStakingTypes.MsgUndelegate{
		DelegatorAddress: "delegator",
		ValidatorAddress: "validator",
		Amount:           cosmosTypes.Coin{Amount: cosmosTypes.NewInt(100000000), Denom: "uatom"},
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

	parsed, err := ParseMsgUndelegate(msgBytes, config.Chains[0], 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	filesystem := &fs.MockFs{}
	logger := loggerPkg.GetNopLogger()
	aliasManager := aliasManagerPkg.NewAliasManager(logger, config, filesystem)
	metricsManager := metrics.NewManager(logger, config.Metrics)
	dataFetcher := data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager)

	err = aliasManager.Set("subscription", "chain", "delegator", "delegator_alias")
	require.NoError(t, err)

	dataFetcher.Cache.Set("chain-id_price_uatom", 6.7)
	dataFetcher.Cache.Set("chain_validator_validator", &responses.Validator{
		OperatorAddress: "test",
		Description:     responses.ValidatorDescription{Moniker: "Validator Moniker"},
	})
	dataFetcher.Cache.Set("chain_staking_params", &responses.StakingParams{
		UnbondingTime: responses.Duration{Duration: 15 * time.Second},
	})

	parsed.GetAdditionalData(dataFetcher, "subscription")

	message, _ := parsed.(*MsgUndelegate)

	require.Equal(t, "delegator_alias", message.DelegatorAddress.Title)
	require.Equal(t, "Validator Moniker", message.ValidatorAddress.Title)
	require.Equal(t, "100.00", fmt.Sprintf("%.2f", message.Amount.Value))
	require.Equal(t, "670.00", fmt.Sprintf("%.2f", message.Amount.PriceUSD))
	require.Equal(t, "atom", message.Amount.Denom.String())
	require.NotZero(t, message.UndelegateFinishTime)
}
