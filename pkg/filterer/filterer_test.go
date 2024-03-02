package filterer_test

import (
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
	filtererPkg "main/pkg/filterer"
	loggerPkg "main/pkg/logger"
	"main/pkg/messages"
	"main/pkg/metrics"
	"main/pkg/types"
	"main/pkg/types/amount"
	"testing"

	queryPkg "github.com/cometbft/cometbft/libs/pubsub/query"
	"github.com/stretchr/testify/require"
)

func TestFilterMessageUnsupported(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{}
	filterer := filtererPkg.NewFilterer(logger, config, nil)

	message := &messages.MsgUnsupportedMessage{}
	require.NotNil(t, filterer.FilterMessage(message, &configTypes.ChainSubscription{
		LogUnknownMessages: true,
		Chain:              "chain",
	}, false))
	require.Nil(t, filterer.FilterMessage(message, &configTypes.ChainSubscription{
		LogUnknownMessages: false,
		Chain:              "chain",
	}, false))
}

func TestFilterMessageUnparsed(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{}
	filterer := filtererPkg.NewFilterer(logger, config, nil)

	message := &messages.MsgUnparsedMessage{}
	require.NotNil(t, filterer.FilterMessage(message, &configTypes.ChainSubscription{
		LogUnparsedMessages: true,
		Chain:               "chain",
	}, false))
	require.Nil(t, filterer.FilterMessage(message, &configTypes.ChainSubscription{
		LogUnparsedMessages: false,
		Chain:               "chain",
	}, false))
}

func TestFilterMessageSimpleNotMatching(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{}
	filterer := filtererPkg.NewFilterer(logger, config, nil)

	message := &messages.MsgSend{
		From: &configTypes.Link{Value: "from"},
		To:   &configTypes.Link{Value: "to"},
		Amount: amount.Amounts{
			amount.AmountFromString("100", "ustake"),
		},
	}
	chainSubscription := &configTypes.ChainSubscription{
		Chain: "chain",
		Filters: configTypes.Filters{
			*queryPkg.MustParse("transfer.sender = 'from2'"),
		},
	}

	require.Nil(t, filterer.FilterMessage(message, chainSubscription, false))
}

func TestFilterMessageSimpleMatching(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{}
	filterer := filtererPkg.NewFilterer(logger, config, nil)

	message := &messages.MsgSend{
		From: &configTypes.Link{Value: "from"},
		To:   &configTypes.Link{Value: "to"},
		Amount: amount.Amounts{
			amount.AmountFromString("100", "ustake"),
		},
	}
	chainSubscription := &configTypes.ChainSubscription{
		Chain: "chain",
		Filters: configTypes.Filters{
			*queryPkg.MustParse("transfer.sender = 'from'"),
		},
	}

	require.NotNil(t, filterer.FilterMessage(message, chainSubscription, false))
}

func TestFilterMessageSimpleRecursiveMatchingExternalAndIgnoreInternal(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{}
	filterer := filtererPkg.NewFilterer(logger, config, nil)

	message := &messages.MsgExec{
		Grantee: &configTypes.Link{Value: "from"},
		Messages: []types.Message{
			&messages.MsgSend{
				From: &configTypes.Link{Value: "from"},
				To:   &configTypes.Link{Value: "to"},
				Amount: amount.Amounts{
					amount.AmountFromString("100", "ustake"),
				},
			},
		},
	}

	chainSubscription := &configTypes.ChainSubscription{
		Chain:                  "chain",
		FilterInternalMessages: false,
		Filters: configTypes.Filters{
			*queryPkg.MustParse("message.action = '/cosmos.authz.v1beta1.MsgExec'"),
		},
	}

	require.NotNil(t, filterer.FilterMessage(message, chainSubscription, false))
}

func TestFilterMessageSimpleRecursiveMatchingExternalAndFilterInternal(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{}
	filterer := filtererPkg.NewFilterer(logger, config, nil)

	message := &messages.MsgExec{
		Grantee: &configTypes.Link{Value: "from"},
		Messages: []types.Message{
			&messages.MsgSend{
				From: &configTypes.Link{Value: "from"},
				To:   &configTypes.Link{Value: "to"},
				Amount: amount.Amounts{
					amount.AmountFromString("100", "ustake"),
				},
			},
		},
	}

	chainSubscription := &configTypes.ChainSubscription{
		Chain:                  "chain",
		FilterInternalMessages: true,
		Filters: configTypes.Filters{
			*queryPkg.MustParse("message.action = '/cosmos.authz.v1beta1.MsgExec'"),
		},
	}

	require.Nil(t, filterer.FilterMessage(message, chainSubscription, false))
}

func TestFilterMessageSimpleRecursiveMatchingExternalAndInternal(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{}
	filterer := filtererPkg.NewFilterer(logger, config, nil)

	message := &messages.MsgExec{
		Grantee: &configTypes.Link{Value: "from"},
		Messages: []types.Message{
			&messages.MsgSend{
				From: &configTypes.Link{Value: "from"},
				To:   &configTypes.Link{Value: "to"},
				Amount: amount.Amounts{
					amount.AmountFromString("100", "ustake"),
				},
			},
		},
	}

	chainSubscription := &configTypes.ChainSubscription{
		Chain:                  "chain",
		FilterInternalMessages: true,
		Filters: configTypes.Filters{
			*queryPkg.MustParse("message.action = '/cosmos.authz.v1beta1.MsgExec'"),
			*queryPkg.MustParse("transfer.sender = 'from'"),
		},
	}

	require.NotNil(t, filterer.FilterMessage(message, chainSubscription, false))
}

func TestFilterReportableTxError(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{}
	logger := loggerPkg.GetDefaultLogger()
	metricsManager := metrics.NewManager(logger, configPkg.MetricsConfig{Enabled: false})
	filterer := filtererPkg.NewFilterer(logger, config, metricsManager)
	chain := &configTypes.Chain{Name: "chain"}

	reportable := &types.TxError{}
	require.NotNil(t, filterer.FilterForChainAndSubscription(reportable, chain, &configTypes.ChainSubscription{
		LogNodeErrors: true,
		Chain:         "chain",
	}))
	require.Nil(t, filterer.FilterForChainAndSubscription(reportable, chain, &configTypes.ChainSubscription{
		LogNodeErrors: false,
		Chain:         "chain",
	}))
}

func TestFilterReportableTxNodeConnectError(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{}
	logger := loggerPkg.GetDefaultLogger()
	metricsManager := metrics.NewManager(logger, configPkg.MetricsConfig{Enabled: false})
	filterer := filtererPkg.NewFilterer(logger, config, metricsManager)
	chain := &configTypes.Chain{Name: "chain"}

	reportable := &types.NodeConnectError{}
	require.NotNil(t, filterer.FilterForChainAndSubscription(reportable, chain, &configTypes.ChainSubscription{
		LogNodeErrors: true,
		Chain:         "chain",
	}))
	require.Nil(t, filterer.FilterForChainAndSubscription(reportable, chain, &configTypes.ChainSubscription{
		LogNodeErrors: false,
		Chain:         "chain",
	}))
}

func TestFilterReportableTxFailed(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{}
	logger := loggerPkg.GetDefaultLogger()
	metricsManager := metrics.NewManager(logger, configPkg.MetricsConfig{Enabled: false})
	filterer := filtererPkg.NewFilterer(logger, config, metricsManager)
	chain := &configTypes.Chain{Name: "chain"}

	reportable := &types.Tx{
		Height: configTypes.Link{Value: "123"},
		Code:   1,
		Messages: []types.Message{
			&messages.MsgSend{
				From: &configTypes.Link{Value: "from"},
				To:   &configTypes.Link{Value: "to"},
				Amount: amount.Amounts{
					amount.AmountFromString("100", "ustake"),
				},
			},
		},
	}
	require.NotNil(t, filterer.FilterForChainAndSubscription(reportable, chain, &configTypes.ChainSubscription{
		LogFailedTransactions: true,
		Chain:                 "chain",
		Filters: configTypes.Filters{
			*queryPkg.MustParse("transfer.sender = 'from'"),
		},
	}))
	require.Nil(t, filterer.FilterForChainAndSubscription(reportable, chain, &configTypes.ChainSubscription{
		LogFailedTransactions: false,
		Chain:                 "chain",
		Filters: configTypes.Filters{
			*queryPkg.MustParse("transfer.sender = 'from'"),
		},
	}))
}

func TestFilterReportableTxProcessedBefore(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{}
	logger := loggerPkg.GetDefaultLogger()
	metricsManager := metrics.NewManager(logger, configPkg.MetricsConfig{Enabled: false})
	filterer := filtererPkg.NewFilterer(logger, config, metricsManager)
	chain := &configTypes.Chain{Name: "chain"}

	subscription := &configTypes.ChainSubscription{
		Chain: "chain",
		Filters: configTypes.Filters{
			*queryPkg.MustParse("transfer.sender = 'from'"),
		},
	}
	reportable := &types.Tx{
		Height: configTypes.Link{Value: "456"},
		Code:   0,
		Messages: []types.Message{
			&messages.MsgSend{
				From: &configTypes.Link{Value: "from"},
				To:   &configTypes.Link{Value: "to"},
				Amount: amount.Amounts{
					amount.AmountFromString("100", "ustake"),
				},
			},
		},
	}
	require.NotNil(t, filterer.FilterForChainAndSubscription(reportable, chain, subscription))

	reportable.Height.Value = "123"

	require.Nil(t, filterer.FilterForChainAndSubscription(reportable, chain, subscription))
}

func TestFilterReportableTxAllMessagesFiltered(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{}
	logger := loggerPkg.GetDefaultLogger()
	metricsManager := metrics.NewManager(logger, configPkg.MetricsConfig{Enabled: false})
	filterer := filtererPkg.NewFilterer(logger, config, metricsManager)
	chain := &configTypes.Chain{Name: "chain"}

	subscription := &configTypes.ChainSubscription{
		Chain: "chain",
		Filters: configTypes.Filters{
			*queryPkg.MustParse("transfer.sender = 'from'"),
		},
	}
	reportable := &types.Tx{
		Height: configTypes.Link{Value: "456"},
		Code:   0,
		Messages: []types.Message{
			&messages.MsgSend{
				From: &configTypes.Link{Value: "from2"},
				To:   &configTypes.Link{Value: "to"},
				Amount: amount.Amounts{
					amount.AmountFromString("100", "ustake"),
				},
			},
		},
	}
	require.Nil(t, filterer.FilterForChainAndSubscription(reportable, chain, subscription))
}

func TestFilterReportableTxInvalidHeight(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	config := &configPkg.AppConfig{}
	logger := loggerPkg.GetDefaultLogger()
	metricsManager := metrics.NewManager(logger, configPkg.MetricsConfig{Enabled: false})
	filterer := filtererPkg.NewFilterer(logger, config, metricsManager)
	chain := &configTypes.Chain{Name: "chain"}

	subscription := &configTypes.ChainSubscription{
		Chain:   "chain",
		Filters: configTypes.Filters{},
	}
	reportable := &types.Tx{
		Height: configTypes.Link{Value: "test"},
		Code:   0,
	}
	_ = filterer.FilterForChainAndSubscription(reportable, chain, subscription)
}

func TestGetReporters(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: configTypes.Chains{
			{Name: "chain-1"},
			{Name: "chain-2"},
		},
		Subscriptions: configTypes.Subscriptions{
			{
				Name:     "subscription-1",
				Reporter: "reporter-1",
				ChainSubscriptions: configTypes.ChainSubscriptions{
					{Chain: "chain-1"},
				},
			},
			{
				Name:     "subscription-2",
				Reporter: "reporter-2",
				ChainSubscriptions: configTypes.ChainSubscriptions{
					{
						Chain:                 "chain-2",
						LogFailedTransactions: false,
					},
					{
						Chain:                 "chain-2",
						LogFailedTransactions: true,
					},
				},
			},
		},
	}
	logger := loggerPkg.GetDefaultLogger()
	metricsManager := metrics.NewManager(logger, configPkg.MetricsConfig{Enabled: false})
	filterer := filtererPkg.NewFilterer(logger, config, metricsManager)
	report := types.Report{
		Chain: &configTypes.Chain{Name: "chain-2"},
		Reportable: &types.Tx{
			Code:   1,
			Height: configTypes.Link{Value: "123"},
			Messages: []types.Message{
				&messages.MsgSend{
					From: &configTypes.Link{Value: "from2"},
					To:   &configTypes.Link{Value: "to"},
					Amount: amount.Amounts{
						amount.AmountFromString("100", "ustake"),
					},
				},
			},
		},
	}
	result := filterer.GetReportableForReporters(report)
	require.Len(t, result, 1)

	_, ok := result["reporter-2"]
	require.True(t, ok)
}
