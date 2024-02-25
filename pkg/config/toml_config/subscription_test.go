package toml_config_test

import (
	tomlConfig "main/pkg/config/toml_config"
	"main/pkg/config/types"
	"testing"

	queryPkg "github.com/cometbft/cometbft/libs/pubsub/query"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestSubscriptionNoName(t *testing.T) {
	t.Parallel()

	subscription := tomlConfig.Subscription{}
	require.Error(t, subscription.Validate())
}

func TestSubscriptionNoReporter(t *testing.T) {
	t.Parallel()

	subscription := tomlConfig.Subscription{
		Name: "name",
	}
	require.Error(t, subscription.Validate())
}

func TestSubscriptionInvalidChainSubscription(t *testing.T) {
	t.Parallel()

	subscription := tomlConfig.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: tomlConfig.ChainSubscriptions{
			{},
		},
	}
	require.Error(t, subscription.Validate())
}

func TestSubscriptionValid(t *testing.T) {
	t.Parallel()

	subscription := tomlConfig.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: tomlConfig.ChainSubscriptions{
			{Chain: "chain"},
		},
	}
	require.NoError(t, subscription.Validate())
}

func TestChainSubscriptionNoName(t *testing.T) {
	t.Parallel()

	subscription := tomlConfig.ChainSubscription{}
	require.Error(t, subscription.Validate())
}

func TestChainSubscriptionInvalidFilter(t *testing.T) {
	t.Parallel()

	subscription := tomlConfig.ChainSubscription{
		Chain:   "chain",
		Filters: []string{"invalid"},
	}
	require.Error(t, subscription.Validate())
}

func TestChainSubscriptionValid(t *testing.T) {
	t.Parallel()

	subscription := tomlConfig.ChainSubscription{
		Chain:   "chain",
		Filters: []string{"event.key = 'value'"},
	}
	require.NoError(t, subscription.Validate())
}

func TestSubscriptionsInvalidSubscription(t *testing.T) {
	t.Parallel()

	subscriptions := tomlConfig.Subscriptions{{}}
	require.Error(t, subscriptions.Validate())
}

func TestSubscriptionsDuplicates(t *testing.T) {
	t.Parallel()

	subscription1 := &tomlConfig.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: tomlConfig.ChainSubscriptions{
			{Chain: "chain"},
		},
	}

	subscription2 := &tomlConfig.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: tomlConfig.ChainSubscriptions{
			{Chain: "chain"},
		},
	}

	subscriptions := tomlConfig.Subscriptions{subscription1, subscription2}
	require.Error(t, subscriptions.Validate())
}

func TestSubscriptionsValid(t *testing.T) {
	t.Parallel()

	subscription := &tomlConfig.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: tomlConfig.ChainSubscriptions{
			{Chain: "chain"},
		},
	}

	subscriptions := tomlConfig.Subscriptions{subscription}
	require.NoError(t, subscriptions.Validate())
}

func TestSubscriptionToAppConfigSubscription(t *testing.T) {
	t.Parallel()

	subscription := &tomlConfig.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: tomlConfig.ChainSubscriptions{
			{Chain: "chain"},
		},
	}
	appConfigSubscription := subscription.ToAppConfigSubscription()

	require.Equal(t, "name", appConfigSubscription.Name)
	require.Equal(t, "reporter", appConfigSubscription.Reporter)
	require.Len(t, appConfigSubscription.ChainSubscriptions, 1)
	require.Equal(t, "chain", appConfigSubscription.ChainSubscriptions[0].Chain)
}

func TestSubscriptionToTomlConfigSubscription(t *testing.T) {
	t.Parallel()

	subscription := &types.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: types.ChainSubscriptions{
			{Chain: "chain"},
		},
	}
	tomlConfigSubscription := tomlConfig.FromAppConfigSubscription(subscription)

	require.Equal(t, "name", tomlConfigSubscription.Name)
	require.Equal(t, "reporter", tomlConfigSubscription.Reporter)
	require.Len(t, tomlConfigSubscription.ChainSubscriptions, 1)
	require.Equal(t, "chain", tomlConfigSubscription.ChainSubscriptions[0].Chain)
}

func TestChainSubscriptionToAppConfigChainSubscription(t *testing.T) {
	t.Parallel()

	subscription := &tomlConfig.ChainSubscription{
		Chain:                  "chain",
		Filters:                []string{"event.key = 'value'"},
		LogUnknownMessages:     null.BoolFrom(true),
		LogUnparsedMessages:    null.BoolFrom(true),
		LogFailedTransactions:  null.BoolFrom(true),
		LogNodeErrors:          null.BoolFrom(true),
		FilterInternalMessages: null.BoolFrom(true),
	}
	appConfigSubscription := subscription.ToAppConfigChainSubscription()

	require.Equal(t, "chain", appConfigSubscription.Chain)
	require.True(t, appConfigSubscription.LogUnknownMessages)
	require.True(t, appConfigSubscription.LogUnparsedMessages)
	require.True(t, appConfigSubscription.LogFailedTransactions)
	require.True(t, appConfigSubscription.LogNodeErrors)
	require.True(t, appConfigSubscription.FilterInternalMessages)
	require.Len(t, appConfigSubscription.Filters, 1)
	require.Equal(t, "event.key = 'value'", appConfigSubscription.Filters[0].String())
}

func TestChainSubscriptionToTomlConfigChainSubscription(t *testing.T) {
	t.Parallel()

	query := queryPkg.MustParse("event.key = 'value'")

	subscription := &types.ChainSubscription{
		Chain:                  "chain",
		Filters:                []queryPkg.Query{*query},
		LogUnknownMessages:     true,
		LogUnparsedMessages:    true,
		LogFailedTransactions:  true,
		LogNodeErrors:          true,
		FilterInternalMessages: true,
	}
	tomlConfigSubscription := tomlConfig.FromAppConfigChainSubscription(subscription)

	require.Equal(t, "chain", tomlConfigSubscription.Chain)
	require.True(t, tomlConfigSubscription.LogUnknownMessages.Bool)
	require.True(t, tomlConfigSubscription.LogUnparsedMessages.Bool)
	require.True(t, tomlConfigSubscription.LogFailedTransactions.Bool)
	require.True(t, tomlConfigSubscription.LogNodeErrors.Bool)
	require.True(t, tomlConfigSubscription.FilterInternalMessages.Bool)
	require.Len(t, tomlConfigSubscription.Filters, 1)
	require.Equal(t, "event.key = 'value'", tomlConfigSubscription.Filters[0])
}
