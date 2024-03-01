package types

import (
	"fmt"
	"strconv"

	"github.com/cometbft/cometbft/libs/pubsub/query"
)

type Chains []*Chain

func (c Chains) FindByName(name string) *Chain {
	for _, chain := range c {
		if chain.Name == name {
			return chain
		}
	}

	return nil
}

func (c Chains) FindByChainID(chainID string) (*Chain, bool) {
	for _, chain := range c {
		if chain.ChainID == chainID {
			return chain, true
		}
	}

	return nil, false
}

func (c Chains) HasChain(name string) bool {
	for _, chain := range c {
		if chain.Name == name {
			return true
		}
	}

	return false
}

type Chain struct {
	Name              string
	PrettyName        string
	ChainID           string
	TendermintNodes   []string
	APINodes          []string
	Queries           []query.Query
	Explorer          *Explorer
	SupportedExplorer SupportedExplorer
	Denoms            DenomInfos
}

func (c Chain) GetName() string {
	if c.PrettyName != "" {
		return c.PrettyName
	}

	return c.Name
}

func (c Chain) GetWalletLink(address string) *Link {
	if c.Explorer == nil {
		return &Link{Value: address}
	}

	return &Link{
		Href:  c.Explorer.GetWalletLink(address),
		Value: address,
	}
}

func (c Chain) GetValidatorLink(address string) *Link {
	if c.Explorer == nil {
		return &Link{Value: address}
	}

	return &Link{
		Href:  fmt.Sprintf(c.Explorer.ValidatorLinkPattern, address),
		Value: address,
	}
}

func (c Chain) GetProposalLink(proposalID string) Link {
	if c.Explorer == nil {
		return Link{Value: proposalID}
	}

	return Link{
		Href:  fmt.Sprintf(c.Explorer.ProposalLinkPattern, proposalID),
		Value: proposalID,
	}
}

func (c Chain) GetTransactionLink(hash string) Link {
	if c.Explorer == nil {
		return Link{Value: hash}
	}

	return Link{
		Href:  fmt.Sprintf(c.Explorer.TransactionLinkPattern, hash),
		Value: hash,
	}
}

func (c Chain) GetBlockLink(height int64) Link {
	heightStr := strconv.FormatInt(height, 10)

	if c.Explorer == nil {
		return Link{Value: heightStr}
	}

	return Link{
		Href:  fmt.Sprintf(c.Explorer.BlockLinkPattern, heightStr),
		Value: heightStr,
	}
}

func (c *Chain) DisplayWarnings() []DisplayWarning {
	var warnings []DisplayWarning

	for _, denom := range c.Denoms {
		warnings = append(warnings, denom.DisplayWarnings(c)...)
	}

	if c.ChainID == "" {
		warnings = append(warnings, DisplayWarning{
			Keys: map[string]string{"chain": c.Name},
			Text: "chain-id is not set, multichain denom matching won't work.",
		})
	}

	if c.Explorer == nil {
		warnings = append(warnings, DisplayWarning{
			Keys: map[string]string{"chain": c.Name},
			Text: "Explorer config not set, links won't be generated.",
		})
	} else {
		warnings = append(warnings, c.Explorer.DisplayWarnings(c)...)
	}

	return warnings
}
