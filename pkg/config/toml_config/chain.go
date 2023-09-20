package toml_config

import (
	"fmt"
	"main/pkg/config/types"

	"github.com/tendermint/tendermint/libs/pubsub/query"
	"gopkg.in/guregu/null.v4"
)

type Chain struct {
	Name                   string     `toml:"name"`
	PrettyName             string     `toml:"pretty-name"`
	TendermintNodes        []string   `toml:"tendermint-nodes"`
	APINodes               []string   `toml:"api-nodes"`
	Queries                []string   `toml:"queries"`
	Filters                []string   `toml:"filters"`
	MintscanPrefix         string     `toml:"mintscan-prefix"`
	PingPrefix             string     `toml:"ping-prefix"`
	PingBaseUrl            string     `default:"https://ping.pub" toml:"ping-base-url"`
	Explorer               *Explorer  `toml:"explorer"`
	LogUnknownMessages     null.Bool  `default:"false"            toml:"log-unknown-messages"`
	LogUnparsedMessages    null.Bool  `default:"true"             toml:"log-unparsed-messages"`
	LogFailedTransactions  null.Bool  `default:"true"             toml:"log-failed-transactions"`
	LogNodeErrors          null.Bool  `default:"true"             toml:"log-node-errors"`
	FilterInternalMessages null.Bool  `default:"true"             toml:"filter-internal-messages"`
	Denoms                 DenomInfos `toml:"denoms"`
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

	for index, denom := range c.Denoms {
		if err := denom.Validate(); err != nil {
			return fmt.Errorf("Error in denom %d: %s", index, err)
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
		Name:                   c.Name,
		PrettyName:             c.PrettyName,
		TendermintNodes:        c.TendermintNodes,
		APINodes:               c.APINodes,
		Queries:                queries,
		Filters:                filters,
		Explorer:               explorer,
		SupportedExplorer:      supportedExplorer,
		LogUnknownMessages:     c.LogUnknownMessages.Bool,
		LogUnparsedMessages:    c.LogUnparsedMessages.Bool,
		LogFailedTransactions:  c.LogFailedTransactions.Bool,
		LogNodeErrors:          c.LogNodeErrors.Bool,
		FilterInternalMessages: c.FilterInternalMessages.Bool,
		Denoms:                 c.Denoms.ToAppConfigDenomInfos(),
	}
}

func FromAppConfigChain(c *types.Chain) *Chain {
	chain := &Chain{
		Name:                   c.Name,
		PrettyName:             c.PrettyName,
		TendermintNodes:        c.TendermintNodes,
		APINodes:               c.APINodes,
		LogUnknownMessages:     null.BoolFrom(c.LogUnknownMessages),
		LogUnparsedMessages:    null.BoolFrom(c.LogUnparsedMessages),
		LogFailedTransactions:  null.BoolFrom(c.LogFailedTransactions),
		LogNodeErrors:          null.BoolFrom(c.LogNodeErrors),
		FilterInternalMessages: null.BoolFrom(c.FilterInternalMessages),
		Denoms:                 TomlConfigDenomsFrom(c.Denoms),
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

type Chains []*Chain
