package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/mcuadros/go-defaults"
	"main/pkg/types/chains"
	"os"
)

type Chains []*chains.Chain

func (c Chains) FindByName(name string) *chains.Chain {
	for _, chain := range c {
		if chain.Name == name {
			return chain
		}
	}

	return nil
}

type Config struct {
	TelegramConfig TelegramConfig `toml:"telegram"`
	LogConfig      LogConfig      `toml:"log"`
	Chains         Chains         `toml:"chains"`
}

type TelegramConfig struct {
	TelegramChat  int64  `toml:"chat"`
	TelegramToken string `toml:"token"`
}

type LogConfig struct {
	LogLevel   string `toml:"level" default:"info"`
	JSONOutput bool   `toml:"json" default:"false"`
}

func (c *Config) Validate() error {
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

func GetConfig(path string) (*Config, error) {
	configBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configString := string(configBytes)

	configStruct := &Config{}
	if _, err = toml.Decode(configString, configStruct); err != nil {
		return nil, err
	}

	defaults.SetDefaults(configStruct)

	for _, chain := range configStruct.Chains {
		if chain.MintscanPrefix != "" {
			chain.Explorer = &chains.Explorer{
				ProposalLinkPattern:    fmt.Sprintf("https://mintscan.io/%s/proposals/%%s", chain.MintscanPrefix),
				WalletLinkPattern:      fmt.Sprintf("https://mintscan.io/%s/account/%%s", chain.MintscanPrefix),
				ValidatorLinkPattern:   fmt.Sprintf("https://mintscan.io/%s/validators/%%s", chain.MintscanPrefix),
				TransactionLinkPattern: fmt.Sprintf("https://mintscan.io/%s/txs/%%s", chain.MintscanPrefix),
				BlockLinkPattern:       fmt.Sprintf("https://mintscan.io/%s/blocks/%%s", chain.MintscanPrefix),
			}
		}

		chain.Filters = make([]chains.Filter, len(chain.FiltersRaw))
		for index, filter := range chain.FiltersRaw {
			chain.Filters[index] = chains.NewFilter(filter)
		}
	}

	return configStruct, nil
}
