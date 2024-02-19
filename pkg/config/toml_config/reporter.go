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
	Name string `toml:"name"`
	Type string `default:"telegram" toml:"type"`

	TelegramConfig *TelegramConfig `toml:"telegram-config"`
}

func (reporter *Reporter) Validate() error {
	if reporter.Name == "" {
		return errors.New("reporter name not provided")
	}

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

type Reporters []*Reporter

func (reporters Reporters) Validate() error {
	for index, reporter := range reporters {
		if err := reporter.Validate(); err != nil {
			return fmt.Errorf("error in reporter %d: %s", index, err)
		}
	}

	// checking names uniqueness
	names := map[string]bool{}

	for _, reporter := range reporters {
		if _, ok := names[reporter.Name]; ok {
			return fmt.Errorf("duplicate reporter name: %s", reporter.Name)
		}

		names[reporter.Name] = true
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
		Name:           reporter.Name,
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
		Name:           reporter.Name,
		Type:           reporter.Type,
		TelegramConfig: telegramConfig,
	}
}
