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
	Timezone      string        `default:"Etc/GMT" toml:"timezone"`

	Reporters []*Reporter `toml:"reporters"`
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

	for index, chain := range c.Chains {
		if err := chain.Validate(); err != nil {
			return fmt.Errorf("error in chain %d: %s", index, err)
		}
	}

	for index, reporter := range c.Reporters {
		if err := reporter.Validate(); err != nil {
			return fmt.Errorf("error in reporter %d: %s", index, err)
		}
	}

	return nil
}
