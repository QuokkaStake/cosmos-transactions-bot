package telegram

import (
	"errors"
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
func TestSetAliasesInvalidInvocationFailedToSend(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	logger := loggerPkg.GetNopLogger()

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "xxx:yyy", Chat: 123, Admins: []int64{1}},
			Timezone:       timezone,
		},
		&configPkg.AppConfig{},
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
			Text:   "/alias",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.Handler(reporter.GetSetAliasCommand())(context)
	require.Error(t, err)
	require.ErrorContains(t, err, "custom error")
}

//nolint:paralleltest // disabled
func TestSetAliasesInvalidInvocation(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Usage: /alias &lt;chain&gt; &lt;address&gt; &lt;alias&gt;"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	logger := loggerPkg.GetNopLogger()

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "xxx:yyy", Chat: 123, Admins: []int64{1}},
			Timezone:       timezone,
		},
		&configPkg.AppConfig{},
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
			Text:   "/alias",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.Handler(reporter.GetSetAliasCommand())(context)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid invocation")
}

//nolint:paralleltest // disabled
func TestSetAliasesDisabled(t *testing.T) {
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

	logger := loggerPkg.GetNopLogger()

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "xxx:yyy", Chat: 123, Admins: []int64{1}},
			Timezone:       timezone,
		},
		&configPkg.AppConfig{},
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
			Text:   "/alias chain address alias",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.Handler(reporter.GetSetAliasCommand())(context)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestSetAliasesCouldNotFindChain(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Chain chain3 is not found in config!"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	logger := loggerPkg.GetNopLogger()
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

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "xxx:yyy", Chat: 123, Admins: []int64{1}},
			Timezone:       timezone,
		},
		&configPkg.AppConfig{},
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
			Text:   "/alias chain3 address alias",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.Handler(reporter.GetSetAliasCommand())(context)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestSetAliasesCouldNotFindSubscription(t *testing.T) {
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

	logger := loggerPkg.GetNopLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "path.yml",
		Chains: configTypes.Chains{
			{Name: "chain1", PrettyName: "Chain1", TendermintNodes: []string{"https://example1.com", "https://example2.com"}},
			{Name: "chain2", PrettyName: "Chain2", TendermintNodes: []string{"https://example3.com", "https://example4.com"}},
		},
		Subscriptions: configTypes.Subscriptions{},
	}
	aliasManager := alias_manager.NewAliasManager(logger, config, &fs.MockFs{})
	metricsManager := metrics.NewManager(logger, configPkg.MetricsConfig{})

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "xxx:yyy", Chat: 123, Admins: []int64{1}},
		},
		config,
		logger,
		nil,
		aliasManager,
		metricsManager,
		data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager),
		"1.2.3",
	)

	err := reporter.Init()
	require.NoError(t, err)

	context := reporter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser"},
			Text:   "/alias chain1 address alias",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.Handler(reporter.GetSetAliasCommand())(context)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestSetAliasesCouldNotSaveAlias(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Error saving alias: yaml: write error: not yet supported"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

	logger := loggerPkg.GetNopLogger()
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
	aliasManager := alias_manager.NewAliasManager(logger, config, &fs.MockFs{FailWrite: true})
	metricsManager := metrics.NewManager(logger, configPkg.MetricsConfig{})

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "xxx:yyy", Chat: 123, Admins: []int64{1}},
		},
		config,
		logger,
		nil,
		aliasManager,
		metricsManager,
		data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager),
		"1.2.3",
	)

	err := reporter.Init()
	require.NoError(t, err)

	context := reporter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser"},
			Text:   "/alias chain1 address alias",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.Handler(reporter.GetSetAliasCommand())(context)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestSetAliasesOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasBytes(assets.GetBytesOrPanic("responses/set-alias.html")),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

	logger := loggerPkg.GetNopLogger()
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

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "xxx:yyy", Chat: 123, Admins: []int64{1}},
		},
		config,
		logger,
		nil,
		aliasManager,
		metricsManager,
		data_fetcher.NewDataFetcher(logger, config, aliasManager, metricsManager),
		"1.2.3",
	)

	err := reporter.Init()
	require.NoError(t, err)

	context := reporter.TelegramBot.NewContext(tele.Update{
		ID: 1,
		Message: &tele.Message{
			Sender: &tele.User{Username: "testuser"},
			Text:   "/alias chain1 address alias",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.Handler(reporter.GetSetAliasCommand())(context)
	require.NoError(t, err)
}
