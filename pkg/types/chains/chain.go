package chains

import (
	"fmt"
	"strconv"
)

type Chain struct {
	Name              string    `toml:"name"`
	PrettyName        string    `toml:"pretty-name"`
	TendermintNodes   []string  `toml:"tendermint-nodes"`
	APINodes          []string  `toml:"api-nodes"`
	Filters           []string  `toml:"filters"`
	MintscanPrefix    string    `toml:"mintscan-prefix"`
	Explorer          *Explorer `toml:"explorer"`
	CoingeckoCurrency string    `toml:"coingecko-currency"`
	BaseDenom         string    `toml:"base-denom"`
	DisplayDenom      string    `toml:"display-denom"`
	DenomCoefficient  int64     `toml:"denom-coefficient" default:"1000000"`
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

	if len(c.Filters) == 0 {
		return fmt.Errorf("no filters provided")
	}

	return nil
}

func (c Chain) GetName() string {
	if c.PrettyName != "" {
		return c.PrettyName
	}

	return c.Name
}

func (c Chain) GetWalletLink(address string) Link {
	if c.Explorer == nil {
		return Link{Title: address}
	}

	return Link{
		Href:  fmt.Sprintf(c.Explorer.WalletLinkPattern, address),
		Title: address,
	}
}

func (c Chain) GetValidatorLink(address string) Link {
	if c.Explorer == nil {
		return Link{Title: address}
	}

	return Link{
		Href:  fmt.Sprintf(c.Explorer.ValidatorLinkPattern, address),
		Title: address,
	}
}

func (c Chain) GetProposalLink(proposalID string) Link {
	if c.Explorer == nil {
		return Link{Title: proposalID}
	}

	return Link{
		Href:  fmt.Sprintf(c.Explorer.ProposalLinkPattern, proposalID),
		Title: proposalID,
	}
}

func (c Chain) GetTransactionLink(hash string) Link {
	if c.Explorer == nil {
		return Link{Title: hash}
	}

	return Link{
		Href:  fmt.Sprintf(c.Explorer.TransactionLinkPattern, hash),
		Title: hash,
	}
}

func (c Chain) GetBlockLink(height int64) Link {
	heightStr := strconv.FormatInt(height, 10)

	if c.Explorer == nil {
		return Link{Title: heightStr}
	}

	return Link{
		Href:  fmt.Sprintf(c.Explorer.BlockLinkPattern, heightStr),
		Title: heightStr,
	}
}
