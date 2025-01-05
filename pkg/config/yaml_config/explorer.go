package yaml_config

import "main/pkg/config/types"

type Explorer struct {
	ProposalLinkPattern    string `yaml:"proposal-link-pattern"`
	WalletLinkPattern      string `yaml:"wallet-link-pattern"`
	ValidatorLinkPattern   string `yaml:"validator-link-pattern"`
	TransactionLinkPattern string `yaml:"transaction-link-pattern"`
	BlockLinkPattern       string `yaml:"block-link-pattern"`
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
