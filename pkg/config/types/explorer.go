package types

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Explorer struct {
	ProposalLinkPattern    string `toml:"proposal-link-pattern"`
	WalletLinkPattern      string `toml:"wallet-link-pattern"`
	ValidatorLinkPattern   string `toml:"validator-link-pattern"`
	TransactionLinkPattern string `toml:"transaction-link-pattern"`
	BlockLinkPattern       string `toml:"block-link-pattern"`
}

func (e *Explorer) DisplayWarnings(logger *zerolog.Logger, c *Chain) {
	if e.ProposalLinkPattern == "" {
		logger.Warn().Str("chain", c.Name).Msg("Proposal link pattern not set, proposals links won't be generated.")
	}

	if e.WalletLinkPattern == "" {
		logger.Warn().Str("chain", c.Name).Msg("Wallet link pattern not set, wallets links won't be generated.")
	}

	if e.ValidatorLinkPattern == "" {
		logger.Warn().Str("chain", c.Name).Msg("Validator link pattern not set, validators links won't be generated.")
	}

	if e.TransactionLinkPattern == "" {
		logger.Warn().Str("chain", c.Name).Msg("Transaction link pattern not set, transactions links won't be generated.")
	}

	if e.BlockLinkPattern == "" {
		logger.Warn().Str("chain", c.Name).Msg("Block link pattern not set, blocks links won't be generated.")
	}
}

type SupportedExplorer interface {
	ToExplorer() *Explorer
}

type MintscanExplorer struct {
	Prefix string
}

func (e *MintscanExplorer) ToExplorer() *Explorer {
	return &Explorer{
		ProposalLinkPattern:    fmt.Sprintf("https://mintscan.io/%s/proposals/%%s", e.Prefix),
		WalletLinkPattern:      fmt.Sprintf("https://mintscan.io/%s/account/%%s", e.Prefix),
		ValidatorLinkPattern:   fmt.Sprintf("https://mintscan.io/%s/validators/%%s", e.Prefix),
		TransactionLinkPattern: fmt.Sprintf("https://mintscan.io/%s/txs/%%s", e.Prefix),
		BlockLinkPattern:       fmt.Sprintf("https://mintscan.io/%s/blocks/%%s", e.Prefix),
	}
}

type PingExplorer struct {
	Prefix  string
	BaseUrl string
}

func (e *PingExplorer) ToExplorer() *Explorer {
	return &Explorer{
		ProposalLinkPattern:    fmt.Sprintf("%s/%s/gov/%%s", e.BaseUrl, e.Prefix),
		WalletLinkPattern:      fmt.Sprintf("%s/%s/account/%%s", e.BaseUrl, e.Prefix),
		ValidatorLinkPattern:   fmt.Sprintf("%s/%s/staking/%%s", e.BaseUrl, e.Prefix),
		TransactionLinkPattern: fmt.Sprintf("%s/%s/tx/%%s", e.BaseUrl, e.Prefix),
		BlockLinkPattern:       fmt.Sprintf("%s/%s/blocks/%%s", e.BaseUrl, e.Prefix),
	}
}
