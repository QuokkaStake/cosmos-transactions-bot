package types

import "fmt"

type Explorer struct {
	ProposalLinkPattern    string `toml:"proposal-link-pattern"`
	WalletLinkPattern      string `toml:"wallet-link-pattern"`
	ValidatorLinkPattern   string `toml:"validator-link-pattern"`
	TransactionLinkPattern string `toml:"transaction-link-pattern"`
	BlockLinkPattern       string `toml:"block-link-pattern"`
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
