package toml_config

import (
	"fmt"
	"main/pkg/config/types"

	"github.com/tendermint/tendermint/libs/pubsub/query"
	null "gopkg.in/guregu/null.v4"
)

type Chain struct {
	Name                  string    `toml:"name"`
	PrettyName            string    `toml:"pretty-name"`
	TendermintNodes       []string  `toml:"tendermint-nodes"`
	APINodes              []string  `toml:"api-nodes"`
	Queries               []string  `toml:"queries"`
	Filters               []string  `toml:"filters"`
	MintscanPrefix        string    `toml:"mintscan-prefix"`
	PingPrefix            string    `toml:"ping-prefix"`
	PingBaseUrl           string    `toml:"ping-base-url" default:"https://ping.pub"`
	Explorer              *Explorer `toml:"explorer"`
	CoingeckoCurrency     string    `toml:"coingecko-currency"`
	BaseDenom             string    `toml:"base-denom"`
	DisplayDenom          string    `toml:"display-denom"`
	DenomCoefficient      int64     `toml:"denom-coefficient" default:"1000000"`
	LogUnknownMessages    null.Bool `toml:"log-unknown-messages" default:"false"`
	LogUnparsedMessages   null.Bool `toml:"log-unparsed-messages" default:"true"`
	LogFailedTransactions null.Bool `toml:"log-failed-transactions" default:"true"`
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

	for index, q := range c.Queries {
		if _, err := query.New(q); err != nil {
			return fmt.Errorf("Error in query %d: %s", index, err)
		}
	}

	for index, filter := range c.Filters {
		if _, err := query.New(filter); err != nil {
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
	} else if c.Explorer != nil {
		explorer = c.Explorer.ToAppConfigExplorer()
	}

	filters := make([]query.Query, len(c.Filters))
	for index, filter := range c.Filters {
		filters[index] = *query.MustParse(filter)
	}

	queries := make([]query.Query, len(c.Queries))
	for index, q := range c.Queries {
		queries[index] = *query.MustParse(q)
	}

	return &types.Chain{
		Name:                  c.Name,
		PrettyName:            c.PrettyName,
		TendermintNodes:       c.TendermintNodes,
		APINodes:              c.APINodes,
		Queries:               queries,
		Filters:               filters,
		Explorer:              explorer,
		SupportedExplorer:     supportedExplorer,
		CoingeckoCurrency:     c.CoingeckoCurrency,
		BaseDenom:             c.BaseDenom,
		DisplayDenom:          c.DisplayDenom,
		DenomCoefficient:      c.DenomCoefficient,
		LogUnknownMessages:    c.LogUnknownMessages.Bool,
		LogUnparsedMessages:   c.LogUnparsedMessages.Bool,
		LogFailedTransactions: c.LogFailedTransactions.Bool,
	}
}

func FromAppConfigChain(c *types.Chain) *Chain {
	chain := &Chain{
		Name:                  c.Name,
		PrettyName:            c.PrettyName,
		TendermintNodes:       c.TendermintNodes,
		APINodes:              c.APINodes,
		CoingeckoCurrency:     c.CoingeckoCurrency,
		BaseDenom:             c.BaseDenom,
		DisplayDenom:          c.DisplayDenom,
		DenomCoefficient:      c.DenomCoefficient,
		LogUnknownMessages:    null.BoolFrom(c.LogUnknownMessages),
		LogUnparsedMessages:   null.BoolFrom(c.LogUnparsedMessages),
		LogFailedTransactions: null.BoolFrom(c.LogFailedTransactions),
	}

	if c.SupportedExplorer == nil && c.Explorer != nil {
		chain.Explorer = &Explorer{
			ProposalLinkPattern:    c.Explorer.ProposalLinkPattern,
			WalletLinkPattern:      c.Explorer.WalletLinkPattern,
			ValidatorLinkPattern:   c.Explorer.ValidatorLinkPattern,
			TransactionLinkPattern: c.Explorer.TransactionLinkPattern,
			BlockLinkPattern:       c.Explorer.BlockLinkPattern,
		}
	} else if mintscan, ok := c.SupportedExplorer.(*types.MintscanExplorer); ok {
		chain.MintscanPrefix = mintscan.Prefix
	} else if ping, ok := c.SupportedExplorer.(*types.PingExplorer); ok {
		chain.PingPrefix = ping.Prefix
		chain.PingBaseUrl = ping.BaseUrl
	}

	chain.Filters = make([]string, len(c.Filters))
	for index, filter := range c.Filters {
		chain.Filters[index] = filter.String()
	}

	chain.Queries = make([]string, len(c.Queries))
	for index, q := range c.Queries {
		chain.Queries[index] = q.String()
	}

	return chain
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
	AliasesPath    string         `toml:"aliases"`
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
	LogLevel   string    `toml:"level" default:"info"`
	JSONOutput null.Bool `toml:"json" default:"false"`
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
