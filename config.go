package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mcuadros/go-defaults"
)

type Explorer struct {
	ProposalLinkPattern  string `toml:"proposal-link-pattern"`
	WalletLinkPattern    string `toml:"wallet-link-pattern"`
	ValidatorLinkPattern string `toml:"validator-link-pattern"`
}

type Chain struct {
	Name            string    `toml:"name"`
	PrettyName      string    `toml:"pretty-name"`
	TendermintNodes []string  `toml:"tendermint-nodes"`
	ApiNodes        []string  `toml:"api-nodes"`
	Filters         []string  `toml:"filters"`
	MintscanPrefix  string    `toml:"mintscan-prefix"`
	Explorer        *Explorer `toml:"explorer"`
}

func (c *Chain) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("empty chain name")
	}

	if len(c.TendermintNodes) == 0 {
		return fmt.Errorf("no Tendermint nodes provided")
	}

	if len(c.ApiNodes) == 0 {
		return fmt.Errorf("no API nodes provided")
	}

	if len(c.Filters) == 0 {
		return fmt.Errorf("no filters provided")
	}

	return nil
}

func (c Chain) GetName() string {
	if c.PrettyName != "" {
		return c.PrettyName
	}

	return c.Name
}

func (c Chain) GetWalletLink(address string) Link {
	if c.Explorer == nil {
		return Link{Title: address}
	}

	return Link{
		Href:  fmt.Sprintf(c.Explorer.WalletLinkPattern, address),
		Title: address,
	}
}

func (c Chain) GetValidatorLink(address string) Link {
	if c.Explorer == nil {
		return Link{Title: address}
	}

	return Link{
		Href:  fmt.Sprintf(c.Explorer.ValidatorLinkPattern, address),
		Title: address,
	}
}

func (c Chain) GetProposalLink(proposalID string) Link {
	if c.Explorer == nil {
		return Link{Title: proposalID}
	}

	return Link{
		Href:  fmt.Sprintf(c.Explorer.ProposalLinkPattern, proposalID),
		Title: proposalID,
	}
}

type Chains []*Chain

func (c Chains) FindByName(name string) *Chain {
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
			chain.Explorer = &Explorer{
				ProposalLinkPattern:  fmt.Sprintf("https://mintscan.io/%s/proposals/%%s", chain.MintscanPrefix),
				WalletLinkPattern:    fmt.Sprintf("https://mintscan.io/%s/account/%%s", chain.MintscanPrefix),
				ValidatorLinkPattern: fmt.Sprintf("https://mintscan.io/%s/validators/%%s", chain.MintscanPrefix),
			}
		}
	}

	return configStruct, nil
}
