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

	cosmosStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestMsgBeginRedelegateParse(t *testing.T) {
	t.Parallel()

	msg := &cosmosStakingTypes.MsgBeginRedelegate{
		DelegatorAddress:    "delegator",
		ValidatorSrcAddress: "validator_src",
		ValidatorDstAddress: "validator_dst",
		Amount:              cosmosTypes.Coin{Amount: cosmosTypes.NewInt(100), Denom: "ustake"},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgBeginRedelegate(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	parsed2, err2 := ParseMsgBeginRedelegate([]byte("aaa"), &configTypes.Chain{Name: "chain"}, 100)
	require.Error(t, err2)
	require.Nil(t, parsed2)
}

func TestMsgBeginRedelegateBase(t *testing.T) {
	t.Parallel()

	msg := &cosmosStakingTypes.MsgBeginRedelegate{
		DelegatorAddress:    "delegator",
		ValidatorSrcAddress: "validator_src",
		ValidatorDstAddress: "validator_dst",
		Amount:              cosmosTypes.Coin{Amount: cosmosTypes.NewInt(100), Denom: "ustake"},
	}
	msgBytes, err := proto.Marshal(msg)
	require.NoError(t, err)

	parsed, err := ParseMsgBeginRedelegate(msgBytes, &configTypes.Chain{Name: "chain"}, 100)
	require.NoError(t, err)
	require.NotNil(t, parsed)

	require.Equal(t, "/cosmos.staking.v1beta1.MsgBeginRedelegate", parsed.Type())

	values := parsed.GetValues()

	require.Equal(t, event.EventValues{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, "/cosmos.staking.v1beta1.MsgBeginRedelegate"),
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeySender, "delegator"),
		event.From(cosmosStakingTypes.EventTypeRedelegate, cosmosStakingTypes.AttributeKeySrcValidator, "validator_src"),
		event.From(cosmosStakingTypes.EventTypeRedelegate, cosmosStakingTypes.AttributeKeyDstValidator, "validator_dst"),
		event.From(cosmosStakingTypes.EventTypeRedelegate, cosmosStakingTypes.AttributeKeyDelegator, "delegator"),
		event.From(cosmosStakingTypes.EventTypeRedelegate, cosmosTypes.AttributeKeyAmount, "100ustake"),
	}, values)

	parsed.AddParsedMessage(nil)
	parsed.SetParsedMessages([]types.Message{})
	require.Empty(t, parsed.GetParsedMessages())
	require.Empty(t, parsed.GetRawMessages())
}

func TestMsgBeginRedelegatePopulate(t *testing.T) {
	t.Parallel()

	msg := &cosmosStakingTypes.MsgBeginRedelegate{
		DelegatorAddress:    "delegator",
		ValidatorSrcAddress: "validator_src",
		ValidatorDstAddress: "validator_dst",
		Amount:              cosmosTypes.Coin{Amount: cosmosTypes.NewInt(100000000), Denom: "uatom"},
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

	parsed, err := ParseMsgBeginRedelegate(msgBytes, config.Chains[0], 100)
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
	dataFetcher.Cache.Set("chain_validator_validator_src", &responses.Validator{
		OperatorAddress: "test",
		Description:     responses.ValidatorDescription{Moniker: "Src Validator Moniker"},
	})
	dataFetcher.Cache.Set("chain_validator_validator_dst", &responses.Validator{
		OperatorAddress: "test",
		Description:     responses.ValidatorDescription{Moniker: "Dst Validator Moniker"},
	})

	parsed.GetAdditionalData(dataFetcher, "subscription")

	msgDelegate, _ := parsed.(*MsgBeginRedelegate)

	require.Equal(t, "delegator_alias", msgDelegate.DelegatorAddress.Title)
	require.Equal(t, "Src Validator Moniker", msgDelegate.ValidatorSrcAddress.Title)
	require.Equal(t, "Dst Validator Moniker", msgDelegate.ValidatorDstAddress.Title)
	require.Equal(t, "100.00", fmt.Sprintf("%.2f", msgDelegate.Amount.Value))
	require.Equal(t, "670.00", fmt.Sprintf("%.2f", msgDelegate.Amount.PriceUSD))
	require.Equal(t, "atom", msgDelegate.Amount.Denom.String())
}
