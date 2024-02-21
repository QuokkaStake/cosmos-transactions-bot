package toml_config

import (
	"fmt"
	"main/pkg/config/types"

	"github.com/cometbft/cometbft/libs/pubsub/query"
)

type Chain struct {
	Name            string     `toml:"name"`
	PrettyName      string     `toml:"pretty-name"`
	TendermintNodes []string   `toml:"tendermint-nodes"`
	APINodes        []string   `toml:"api-nodes"`
	Queries         []string   `default:"[\"tx.height > 1\"]" toml:"queries"`
	MintscanPrefix  string     `toml:"mintscan-prefix"`
	PingPrefix      string     `toml:"ping-prefix"`
	PingBaseUrl     string     `default:"https://ping.pub"    toml:"ping-base-url"`
	Explorer        *Explorer  `toml:"explorer"`
	Denoms          DenomInfos `toml:"denoms"`
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
			return fmt.Errorf("error in query %d: %s", index, err)
		}
	}

	for index, denom := range c.Denoms {
		if err := denom.Validate(); err != nil {
			return fmt.Errorf("error in denom %d: %s", index, err)
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

	queries := make([]query.Query, len(c.Queries))
	for index, q := range c.Queries {
		queries[index] = *query.MustParse(q)
	}

	return &types.Chain{
		Name:              c.Name,
		PrettyName:        c.PrettyName,
		TendermintNodes:   c.TendermintNodes,
		APINodes:          c.APINodes,
		Queries:           queries,
		Explorer:          explorer,
		SupportedExplorer: supportedExplorer,
		Denoms:            c.Denoms.ToAppConfigDenomInfos(),
	}
}

func FromAppConfigChain(c *types.Chain) *Chain {
	chain := &Chain{
		Name:            c.Name,
		PrettyName:      c.PrettyName,
		TendermintNodes: c.TendermintNodes,
		APINodes:        c.APINodes,
		Denoms:          TomlConfigDenomsFrom(c.Denoms),
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

	chain.Queries = make([]string, len(c.Queries))
	for index, q := range c.Queries {
		chain.Queries[index] = q.String()
	}

	return chain
}

type Chains []*Chain

func (chains Chains) Validate() error {
	for index, chain := range chains {
		if err := chain.Validate(); err != nil {
			return fmt.Errorf("error in chain %d: %s", index, err)
		}
	}

	// checking names uniqueness
	names := map[string]bool{}

	for _, chain := range chains {
		if _, ok := names[chain.Name]; ok {
			return fmt.Errorf("duplicate chain name: %s", chain.Name)
		}

		names[chain.Name] = true
	}

	return nil
}

func (chains Chains) HasChainByName(name string) bool {
	for _, chain := range chains {
		if chain.Name == name {
			return true
		}
	}

	return false
}
