package yaml_config

import (
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
)

type YamlConfig struct {
	AliasesPath   string        `yaml:"aliases"`
	LogConfig     LogConfig     `yaml:"log"`
	MetricsConfig MetricsConfig `yaml:"metrics"`
	Chains        Chains        `yaml:"chains"`
	Subscriptions Subscriptions `yaml:"subscriptions"`
	Timezone      string        `default:"Etc/GMT"    yaml:"timezone"`

	Reporters Reporters `yaml:"reporters"`
}

type LogConfig struct {
	LogLevel   string    `default:"info"  yaml:"level"`
	JSONOutput null.Bool `default:"false" yaml:"json"`
}

func (c *YamlConfig) Validate() error {
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
		for chainSubscriptionIndex, chainSubscription := range subscription.ChainSubscriptions {
			if !c.Chains.HasChainByName(chainSubscription.Chain) {
				return fmt.Errorf(
					"error in subscription %d: error in chain %d: no such chain '%s'",
					index,
					chainSubscriptionIndex,
					chainSubscription.Chain,
				)
			}
		}

		if !c.Reporters.HasReporterByName(subscription.Reporter) {
			return fmt.Errorf(
				"error in subscription %d: no such reporter '%s'",
				index,
				subscription.Reporter,
			)
		}
	}

	return nil
}
