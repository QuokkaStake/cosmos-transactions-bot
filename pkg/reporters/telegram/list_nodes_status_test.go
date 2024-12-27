package telegram

import (
	"main/assets"
	"main/pkg/alias_manager"
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
	"main/pkg/data_fetcher"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"
	"main/pkg/metrics"
	"main/pkg/nodes_manager"
	"main/pkg/types"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	tele "gopkg.in/telebot.v3"
)

//nolint:paralleltest // disabled
func TestListNodesCouldNotFindSubscription(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("This reporter is not linked to any chains!"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{AliasesPath: "path.yml"}
	aliasManager := alias_manager.NewAliasManager(logger, config, &fs.MockFs{})
	metricsManager := metrics.NewManager(logger, configPkg.MetricsConfig{})

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "xxx:yyy", Chat: 123, Admins: []int64{1}},
		},
		&configPkg.AppConfig{Timezone: timezone},
		logger,
		nil,
		aliasManager,
		metricsManager,
		data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager),
		"1.2.3",
	)

	err = reporter.Init()
	require.NoError(t, err)

	context := reporter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser"},
			Text:   "/aliases",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.Handler(reporter.GetListNodesCommand())(context)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestListNodesOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasBytes(assets.GetBytesOrPanic("responses/status.html")),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "path.yml",
		Chains: configTypes.Chains{
			{Name: "chain1", PrettyName: "Chain1", TendermintNodes: []string{"https://example1.com", "https://example2.com"}},
			{Name: "chain2", PrettyName: "Chain2", TendermintNodes: []string{"https://example3.com", "https://example4.com"}},
		},
		Subscriptions: configTypes.Subscriptions{
			{
				Name:     "subscription",
				Reporter: "reporter",
				ChainSubscriptions: configTypes.ChainSubscriptions{{
					Chain: "chain1",
				}},
			},
		},
	}
	aliasManager := alias_manager.NewAliasManager(logger, config, &fs.MockFs{})
	metricsManager := metrics.NewManager(logger, configPkg.MetricsConfig{})
	nodeManager := nodes_manager.NewNodesManager(logger, config, metricsManager)

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "xxx:yyy", Chat: 123, Admins: []int64{1}},
		},
		&configPkg.AppConfig{Timezone: timezone},
		logger,
		nodeManager,
		aliasManager,
		metricsManager,
		data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager),
		"1.2.3",
	)

	err = reporter.Init()
	require.NoError(t, err)

	context := reporter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser"},
			Text:   "/aliases",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.Handler(reporter.GetListNodesCommand())(context)
	require.NoError(t, err)
}
