package telegram

import (
	"main/assets"
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
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
func TestAppHelpOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/getMe",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("telegram-bot-ok.json")))

	httpmock.RegisterMatcherResponder(
		"POST",
		"https://api.telegram.org/botxxx:yyy/sendMessage",
		types.TelegramResponseHasBytes(assets.GetBytesOrPanic("responses/help.html")),
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
		nil,
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
			Text:   "/help",
			Chat:   &tele.Chat{ID: 2},
		},
	})

	err = reporter.Handler(reporter.GetHelpCommand())(context)
	require.NoError(t, err)
}
