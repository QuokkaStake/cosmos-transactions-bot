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

	tele "gopkg.in/telebot.v3"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // disabled
func TestTelegramReporterNoCredentials(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewErrorResponder(errors.New("custom error")))

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "", Chat: 0},
		},
		&configPkg.AppConfig{Timezone: timezone},
		loggerPkg.GetDefaultLogger(),
		nil,
		nil,
		nil,
		nil,
		"v1.2.3",
	)

	err = reporter.Init()
	require.NoError(t, err)
}

//nolint:paralleltest // disabled
func TestTelegramReporterCannotFetchBot(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewErrorResponder(errors.New("custom error")))

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "xxx:yyy", Chat: 123},
		},
		&configPkg.AppConfig{Timezone: timezone},
		loggerPkg.GetDefaultLogger(),
		nil,
		nil,
		nil,
		nil,
		"v1.2.3",
	)

	err = reporter.Init()
	require.Error(t, err)
	require.ErrorContains(t, err, "custom error")
}

//nolint:paralleltest // disabled
func TestTelegramReporterStartsOkay(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "xxx:yyy", Chat: 123, Admins: []int64{1}},
		},
		&configPkg.AppConfig{Timezone: timezone},
		loggerPkg.GetDefaultLogger(),
		nil,
		nil,
		nil,
		nil,
		"v1.2.3",
	)

	err = reporter.Init()
	require.NoError(t, err)

	go reporter.Start()
	reporter.Stop()
}

//nolint:paralleltest // disabled
func TestTelegramReporterStartsDisabled(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	reporter := NewReporter(
		&configTypes.Reporter{
			Name:           "reporter",
			Type:           "telegram",
			TelegramConfig: &configTypes.TelegramConfig{Token: "xxx:yyy", Chat: 123, Admins: []int64{1}},
		},
		&configPkg.AppConfig{Timezone: timezone},
		loggerPkg.GetDefaultLogger(),
		nil,
		nil,
		nil,
		nil,
		"v1.2.3",
	)

	reporter.Start()
}

//nolint:paralleltest // disabled
func TestHandlerInternalError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasText("Internal error!"),
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-send-message-ok.json")),
	)

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{}
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
			Text:   "/command",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	command := Command{
		Name:    "command",
		MinArgs: 0,
		Usage:   "/command",
		Query:   "command",
		Execute: func(c tele.Context) (string, error) {
			return "", errors.New("custom error")
		},
	}

	err = reporter.Handler(command)(context)
	require.NoError(t, err)
}
