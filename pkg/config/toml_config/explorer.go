package toml_config

import "main/pkg/config/types"

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
