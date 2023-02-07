package types

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/tendermint/tendermint/libs/pubsub/query"
)

type Chain struct {
	Name              string
	PrettyName        string
	TendermintNodes   []string
	APINodes          []string
	Queries           []query.Query
	Explorer          *Explorer
	SupportedExplorer SupportedExplorer
	Denoms            DenomInfos

	LogUnknownMessages     bool
	LogUnparsedMessages    bool
	LogFailedTransactions  bool
	FilterInternalMessages bool

	Filters Filters
}

func (c Chain) GetName() string {
	if c.PrettyName != "" {
		return c.PrettyName
	}

	return c.Name
}

func (c Chain) GetWalletLink(address string) Link {
	if c.Explorer == nil {
		return Link{Value: address}
	}

	return Link{
		Href:  fmt.Sprintf(c.Explorer.WalletLinkPattern, address),
		Value: address,
	}
}

func (c Chain) GetValidatorLink(address string) Link {
	if c.Explorer == nil {
		return Link{Value: address}
	}

	return Link{
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

func (c *Chain) DisplayWarnings(logger *zerolog.Logger) {
	if len(c.Denoms) == 0 {
		logger.Warn().Str("chain", c.Name).Msg("No denoms set, prices in USD won't be displayed.")
	} else {
		for _, denom := range c.Denoms {
			denom.DisplayWarnings(c, logger)
		}
	}

	if c.Explorer == nil {
		logger.Warn().Str("chain", c.Name).Msg("Explorer config not set, links won't be generated.")
	} else {
		c.Explorer.DisplayWarnings(logger, c)
	}
}
