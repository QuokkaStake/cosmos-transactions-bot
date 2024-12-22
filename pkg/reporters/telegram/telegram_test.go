package telegram

import (
	"errors"
	"main/assets"
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
	loggerPkg "main/pkg/logger"
	"testing"
	"time"

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
