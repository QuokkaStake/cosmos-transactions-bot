package toml_config

import (
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
)

type TomlConfig struct {
	AliasesPath   string        `toml:"aliases"`
	LogConfig     LogConfig     `toml:"log"`
	MetricsConfig MetricsConfig `toml:"metrics"`
	Chains        Chains        `toml:"chains"`
	Subscriptions Subscriptions `toml:"subscriptions"`
	Timezone      string        `default:"Etc/GMT" toml:"timezone"`

	Reporters Reporters `toml:"reporters"`
}

type LogConfig struct {
	LogLevel   string    `default:"info"  toml:"level"`
	JSONOutput null.Bool `default:"false" toml:"json"`
}

func (c *TomlConfig) Validate() error {
	if len(c.Chains) == 0 {
		return fmt.Errorf("no chains provided")
	}

	if _, err := time.LoadLocation(c.Timezone); err != nil {
		return fmt.Errorf("error parsing timezone: %s", err)
	}

	if err := c.Chains.Validate(); err != nil {
		return fmt.Errorf("error in chains: %s", err)
	}

	if err := c.Reporters.Validate(); err != nil {
		return fmt.Errorf("error in reporters: %s", err)
	}

	if err := c.Subscriptions.Validate(); err != nil {
		return fmt.Errorf("error in subscriptions: %s", err)
	}

	for index, subscription := range c.Subscriptions {
		if !c.Chains.HasChainByName(subscription.Chain) {
			return fmt.Errorf("error in subscription %d: no such chain '%s'", index, subscription.Chain)
		}

		if !c.Reporters.HasReporterByName(subscription.Reporter) {
			return fmt.Errorf("error in subscription %d: no such chain '%s'", index, subscription.Chain)
		}
	}

	return nil
}
