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
	"main/pkg/types"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	tele "gopkg.in/telebot.v3"
)

//nolint:paralleltest // disabled
func TestGetAliasesDisabled(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Aliases manager is not enabled!"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	logger := loggerPkg.GetDefaultLogger()

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "xxx:yyy", Chat: 123, Admins: []int64{1}},
		},
		&configPkg.AppConfig{Timezone: timezone},
		logger,
		nil,
		alias_manager.NewAliasManager(logger, &configPkg.AppConfig{}, &fs.MockFs{}),
		metrics.NewManager(logger, configPkg.MetricsConfig{}),
		nil,
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

	err = reporter.Handler(reporter.GetGetAliasesCommand())(context)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestGetAliasesCouldNotFindSubscription(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("This reporter is not linked to any subscription!"),
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

	err = reporter.Handler(reporter.GetGetAliasesCommand())(context)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestGetAliasesOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasBytes(assets.GetBytesOrPanic("responses/get-aliases.html")),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "aliases.toml",
		Chains:      configTypes.Chains{{Name: "chain", PrettyName: "ChainName"}},
		Subscriptions: configTypes.Subscriptions{{
			Name:     "subscription",
			Reporter: "reporter",
			ChainSubscriptions: configTypes.ChainSubscriptions{{
				Chain: "chain",
			}},
		}},
	}
	aliasManager := alias_manager.NewAliasManager(logger, config, &fs.MockFs{})
	metricsManager := metrics.NewManager(logger, configPkg.MetricsConfig{})

	aliasManager.Load()

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

	err = reporter.Handler(reporter.GetGetAliasesCommand())(context)
	require.NoError(t, err)
}
