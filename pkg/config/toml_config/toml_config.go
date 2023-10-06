package toml_config

import (
	"fmt"

	"gopkg.in/guregu/null.v4"
)

type TomlConfig struct {
	AliasesPath    string         `toml:"aliases"`
	TelegramConfig TelegramConfig `toml:"telegram"`
	LogConfig      LogConfig      `toml:"log"`
	MetricsConfig  MetricsConfig  `toml:"metrics"`
	Chains         Chains         `toml:"chains"`
}

type TelegramConfig struct {
	Chat   int64   `toml:"chat"`
	Token  string  `toml:"token"`
	Admins []int64 `toml:"admins"`
}

type LogConfig struct {
	LogLevel   string    `default:"info"  toml:"level"`
	JSONOutput null.Bool `default:"false" toml:"json"`
}

func (c *TomlConfig) Validate() error {
	if len(c.Chains) == 0 {
		return fmt.Errorf("no chains provided")
	}

	for index, chain := range c.Chains {
		if err := chain.Validate(); err != nil {
			return fmt.Errorf("error in chain %d: %s", index, err)
		}
	}

	return nil
}
