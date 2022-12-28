package types

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog"
)

type Chain struct {
	Name               string
	PrettyName         string
	TendermintNodes    []string
	APINodes           []string
	Queries            []string
	Explorer           *Explorer
	SupportedExplorer  SupportedExplorer
	CoingeckoCurrency  string
	BaseDenom          string
	DisplayDenom       string
	DenomCoefficient   int64
	LogUnknownMessages bool

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
		Title: hash,
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

func (c Chain) DisplayWarnings(log *zerolog.Logger) {
	if c.BaseDenom == "" {
		log.Warn().Str("chain", c.Name).Msg("Base denom not set, denoms won't be displayed correctly.")
	}

	if c.DisplayDenom == "" {
		log.Warn().Str("chain", c.Name).Msg("Display denom not set, denoms won't be displayed correctly.")
	}

	if c.CoingeckoCurrency == "" {
		log.Warn().Str("chain", c.Name).Msg("Coingecko currency not set, prices in USD won't be displayed.")
	}

	if c.Explorer == nil {
		log.Warn().Str("chain", c.Name).Msg("Explorer config not set, links won't be generated.")
	}
}
