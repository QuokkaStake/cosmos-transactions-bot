package toml_config

import (
	"fmt"
	"strings"

	"main/pkg/config/types"
	"main/pkg/utils"
)

type Chain struct {
	Name               string    `toml:"name"`
	PrettyName         string    `toml:"pretty-name"`
	TendermintNodes    []string  `toml:"tendermint-nodes"`
	APINodes           []string  `toml:"api-nodes"`
	Queries            []string  `toml:"queries"`
	Filters            []string  `toml:"filters"`
	MintscanPrefix     string    `toml:"mintscan-prefix"`
	PingPrefix         string    `toml:"ping-prefix"`
	PingBaseUrl        string    `toml:"ping-base-url" default:"https://ping.pub"`
	Explorer           *Explorer `toml:"explorer" default:"-"`
	CoingeckoCurrency  string    `toml:"coingecko-currency"`
	BaseDenom          string    `toml:"base-denom"`
	DisplayDenom       string    `toml:"display-denom"`
	DenomCoefficient   int64     `toml:"denom-coefficient" default:"1000000"`
	LogUnknownMessages bool      `toml:"log-unknown-messages" default:"false"`
}

func (c *Chain) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("empty chain name")
	}

	if len(c.TendermintNodes) == 0 {
		return fmt.Errorf("no Tendermint nodes provided")
	}

	if len(c.APINodes) == 0 {
		return fmt.Errorf("no API nodes provided")
	}

	if len(c.Queries) == 0 {
		return fmt.Errorf("no queries provided")
	}

	for index, filter := range c.Filters {
		if err := ValidateFilter(filter); err != nil {
			return fmt.Errorf("Error in filter %d: %s", index, err)
		}
	}

	return nil
}

func (c *Chain) ToAppConfigChain() *types.Chain {
	var supportedExplorer types.SupportedExplorer

	if c.MintscanPrefix != "" {
		supportedExplorer = &types.MintscanExplorer{Prefix: c.MintscanPrefix}
	} else if c.PingPrefix != "" {
		supportedExplorer = &types.PingExplorer{Prefix: c.PingPrefix, BaseUrl: c.PingBaseUrl}
	}

	var explorer *types.Explorer
	if supportedExplorer != nil {
		explorer = supportedExplorer.ToExplorer()
	} else if explorer != nil {
		explorer = c.Explorer.ToAppConfigExplorer()
	}

	filters := make([]types.Filter, len(c.Filters))
	for index, filter := range c.Filters {
		filters[index] = types.NewFilter(filter)
	}

	return &types.Chain{
		Name:               c.Name,
		PrettyName:         c.PrettyName,
		TendermintNodes:    c.TendermintNodes,
		APINodes:           c.APINodes,
		Queries:            c.Queries,
		Filters:            filters,
		Explorer:           explorer,
		SupportedExplorer:  supportedExplorer,
		CoingeckoCurrency:  c.CoingeckoCurrency,
		BaseDenom:          c.BaseDenom,
		DisplayDenom:       c.DisplayDenom,
		DenomCoefficient:   c.DenomCoefficient,
		LogUnknownMessages: c.LogUnknownMessages,
	}
}

func ValidateFilter(filter string) error {
	split := strings.Split(filter, " ")
	if len(split) != 3 {
		return fmt.Errorf("filter should match pattern: <key> <operator> <value>")
	}

	if !utils.Contains([]string{"=", "!="}, split[1]) {
		return fmt.Errorf("unknown operator %s, allowed are: '=', '!='", split[1])
	}

	return nil
}

type Explorer struct {
	ProposalLinkPattern    string `toml:"proposal-link-pattern"`
	WalletLinkPattern      string `toml:"wallet-link-pattern"`
	ValidatorLinkPattern   string `toml:"validator-link-pattern"`
	TransactionLinkPattern string `toml:"transaction-link-pattern"`
	BlockLinkPattern       string `toml:"block-link-pattern"`
}

func (e *Explorer) ToAppConfigExplorer() *types.Explorer {
	return &types.Explorer{
		ProposalLinkPattern:    e.ProposalLinkPattern,
		WalletLinkPattern:      e.WalletLinkPattern,
		ValidatorLinkPattern:   e.ValidatorLinkPattern,
		TransactionLinkPattern: e.TransactionLinkPattern,
		BlockLinkPattern:       e.BlockLinkPattern,
	}
}

type Chains []*Chain

type TomlConfig struct {
	TelegramConfig TelegramConfig `toml:"telegram"`
	LogConfig      LogConfig      `toml:"log"`
	Chains         Chains         `toml:"chains"`
}

type TelegramConfig struct {
	Chat   int64   `toml:"chat"`
	Token  string  `toml:"token"`
	Admins []int64 `toml:"admins"`
}

type LogConfig struct {
	LogLevel   string `toml:"level" default:"info"`
	JSONOutput bool   `toml:"json" default:"false"`
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
