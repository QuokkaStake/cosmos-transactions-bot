package types

import (
	"fmt"
)

type Explorer struct {
	ProposalLinkPattern    string `toml:"proposal-link-pattern"`
	WalletLinkPattern      string `toml:"wallet-link-pattern"`
	ValidatorLinkPattern   string `toml:"validator-link-pattern"`
	TransactionLinkPattern string `toml:"transaction-link-pattern"`
	BlockLinkPattern       string `toml:"block-link-pattern"`
}

func (e *Explorer) GetWalletLink(address string) string {
	return fmt.Sprintf(e.WalletLinkPattern, address)
}

func (e *Explorer) DisplayWarnings(c *Chain) []DisplayWarning {
	var warnings []DisplayWarning

	if e.ProposalLinkPattern == "" {
		warnings = append(warnings, DisplayWarning{
			Keys: map[string]string{
				"chain": c.Name,
			},
			Text: "Proposal link pattern not set, proposals links won't be generated.",
		})
	}

	if e.WalletLinkPattern == "" {
		warnings = append(warnings, DisplayWarning{
			Keys: map[string]string{
				"chain": c.Name,
			},
			Text: "Wallet link pattern not set, wallets links won't be generated.",
		})
	}

	if e.ValidatorLinkPattern == "" {
		warnings = append(warnings, DisplayWarning{
			Keys: map[string]string{
				"chain": c.Name,
			},
			Text: "Validator link pattern not set, validators links won't be generated.",
		})
	}

	if e.TransactionLinkPattern == "" {
		warnings = append(warnings, DisplayWarning{
			Keys: map[string]string{
				"chain": c.Name,
			},
			Text: "Transaction link pattern not set, transactions links won't be generated.",
		})
	}

	if e.BlockLinkPattern == "" {
		warnings = append(warnings, DisplayWarning{
			Keys: map[string]string{
				"chain": c.Name,
			},
			Text: "Block link pattern not set, blocks links won't be generated.",
		})
	}

	return warnings
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
		TransactionLinkPattern: fmt.Sprintf("https://mintscan.io/%s/tx/%%s", e.Prefix),
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
