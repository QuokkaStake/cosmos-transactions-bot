package toml_config

import (
	"errors"
	"fmt"
	"main/pkg/config/types"
	"main/pkg/constants"
	"main/pkg/utils"
	"strings"
)

type TelegramConfig struct {
	Chat   int64   `toml:"chat"`
	Token  string  `toml:"token"`
	Admins []int64 `toml:"admins"`
}

type Reporter struct {
	Type string `default:"telegram" toml:"type"`

	TelegramConfig *TelegramConfig `toml:"telegram-config"`
}

func (reporter *Reporter) Validate() error {
	reporterTypes := constants.GetReporterTypes()
	if !utils.Contains(reporterTypes, reporter.Type) {
		return fmt.Errorf(
			"expected type to be one of %s, but got %s",
			strings.Join(reporterTypes, ", "),
			reporter.Type,
		)
	}

	if reporter.Type == constants.ReporterTypeTelegram && reporter.TelegramConfig == nil {
		return errors.New("missing telegram-config for Telegram reporter")
	}

	return nil
}

func FromAppConfigReporter(reporter *types.Reporter) *Reporter {
	var telegramConfig *TelegramConfig

	if reporter.TelegramConfig != nil {
		telegramConfig = &TelegramConfig{
			Chat:   reporter.TelegramConfig.Chat,
			Token:  reporter.TelegramConfig.Token,
			Admins: reporter.TelegramConfig.Admins,
		}
	}

	return &Reporter{
		Type:           reporter.Type,
		TelegramConfig: telegramConfig,
	}
}

func (reporter *Reporter) ToAppConfigReporter() *types.Reporter {
	var telegramConfig *types.TelegramConfig

	if reporter.TelegramConfig != nil {
		telegramConfig = &types.TelegramConfig{
			Chat:   reporter.TelegramConfig.Chat,
			Token:  reporter.TelegramConfig.Token,
			Admins: reporter.TelegramConfig.Admins,
		}
	}

	return &types.Reporter{
		Type:           reporter.Type,
		TelegramConfig: telegramConfig,
	}
}
