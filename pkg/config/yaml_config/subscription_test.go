package yaml_config_test

import (
	"main/pkg/config/types"
	yamlConfig "main/pkg/config/yaml_config"
	"testing"

	queryPkg "github.com/cometbft/cometbft/libs/pubsub/query"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func TestSubscriptionNoName(t *testing.T) {
	t.Parallel()

	subscription := yamlConfig.Subscription{}
	require.Error(t, subscription.Validate())
}

func TestSubscriptionNoReporter(t *testing.T) {
	t.Parallel()

	subscription := yamlConfig.Subscription{
		Name: "name",
	}
	require.Error(t, subscription.Validate())
}

func TestSubscriptionInvalidChainSubscription(t *testing.T) {
	t.Parallel()

	subscription := yamlConfig.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: yamlConfig.ChainSubscriptions{
			{},
		},
	}
	require.Error(t, subscription.Validate())
}

func TestSubscriptionValid(t *testing.T) {
	t.Parallel()

	subscription := yamlConfig.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: yamlConfig.ChainSubscriptions{
			{Chain: "chain"},
		},
	}
	require.NoError(t, subscription.Validate())
}

func TestChainSubscriptionNoName(t *testing.T) {
	t.Parallel()

	subscription := yamlConfig.ChainSubscription{}
	require.Error(t, subscription.Validate())
}

func TestChainSubscriptionInvalidFilter(t *testing.T) {
	t.Parallel()

	subscription := yamlConfig.ChainSubscription{
		Chain:   "chain",
		Filters: []string{"invalid"},
	}
	require.Error(t, subscription.Validate())
}

func TestChainSubscriptionValid(t *testing.T) {
	t.Parallel()

	subscription := yamlConfig.ChainSubscription{
		Chain:   "chain",
		Filters: []string{"event.key = 'value'"},
	}
	require.NoError(t, subscription.Validate())
}

func TestSubscriptionsInvalidSubscription(t *testing.T) {
	t.Parallel()

	subscriptions := yamlConfig.Subscriptions{{}}
	require.Error(t, subscriptions.Validate())
}

func TestSubscriptionsDuplicates(t *testing.T) {
	t.Parallel()

	subscription1 := &yamlConfig.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: yamlConfig.ChainSubscriptions{
			{Chain: "chain"},
		},
	}

	subscription2 := &yamlConfig.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: yamlConfig.ChainSubscriptions{
			{Chain: "chain"},
		},
	}

	subscriptions := yamlConfig.Subscriptions{subscription1, subscription2}
	require.Error(t, subscriptions.Validate())
}

func TestSubscriptionsValid(t *testing.T) {
	t.Parallel()

	subscription := &yamlConfig.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: yamlConfig.ChainSubscriptions{
			{Chain: "chain"},
		},
	}

	subscriptions := yamlConfig.Subscriptions{subscription}
	require.NoError(t, subscriptions.Validate())
}

func TestSubscriptionToAppConfigSubscription(t *testing.T) {
	t.Parallel()

	subscription := &yamlConfig.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: yamlConfig.ChainSubscriptions{
			{Chain: "chain"},
		},
	}
	appConfigSubscription := subscription.ToAppConfigSubscription()

	require.Equal(t, "name", appConfigSubscription.Name)
	require.Equal(t, "reporter", appConfigSubscription.Reporter)
	require.Len(t, appConfigSubscription.ChainSubscriptions, 1)
	require.Equal(t, "chain", appConfigSubscription.ChainSubscriptions[0].Chain)
}

func TestSubscriptionToYamlConfigSubscription(t *testing.T) {
	t.Parallel()

	subscription := &types.Subscription{
		Name:     "name",
		Reporter: "reporter",
		ChainSubscriptions: types.ChainSubscriptions{
			{Chain: "chain"},
		},
	}
	yamlConfigSubscription := yamlConfig.FromAppConfigSubscription(subscription)

	require.Equal(t, "name", yamlConfigSubscription.Name)
	require.Equal(t, "reporter", yamlConfigSubscription.Reporter)
	require.Len(t, yamlConfigSubscription.ChainSubscriptions, 1)
	require.Equal(t, "chain", yamlConfigSubscription.ChainSubscriptions[0].Chain)
}

func TestChainSubscriptionToAppConfigChainSubscription(t *testing.T) {
	t.Parallel()

	subscription := &yamlConfig.ChainSubscription{
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

func TestChainSubscriptionToYamlConfigChainSubscription(t *testing.T) {
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
	yamlConfigSubscription := yamlConfig.FromAppConfigChainSubscription(subscription)

	require.Equal(t, "chain", yamlConfigSubscription.Chain)
	require.True(t, yamlConfigSubscription.LogUnknownMessages.Bool)
	require.True(t, yamlConfigSubscription.LogUnparsedMessages.Bool)
	require.True(t, yamlConfigSubscription.LogFailedTransactions.Bool)
	require.True(t, yamlConfigSubscription.LogNodeErrors.Bool)
	require.True(t, yamlConfigSubscription.FilterInternalMessages.Bool)
	require.Len(t, yamlConfigSubscription.Filters, 1)
	require.Equal(t, "event.key = 'value'", yamlConfigSubscription.Filters[0])
}
