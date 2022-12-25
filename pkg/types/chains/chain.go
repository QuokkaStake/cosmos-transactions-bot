package chains

import (
	"fmt"
	"main/pkg/utils"
	"strconv"
	"strings"
)

type Chain struct {
	Name               string    `toml:"name"`
	PrettyName         string    `toml:"pretty-name"`
	TendermintNodes    []string  `toml:"tendermint-nodes"`
	APINodes           []string  `toml:"api-nodes"`
	Queries            []string  `toml:"queries"`
	FiltersRaw         []string  `toml:"filters"`
	MintscanPrefix     string    `toml:"mintscan-prefix"`
	Explorer           *Explorer `toml:"explorer"`
	CoingeckoCurrency  string    `toml:"coingecko-currency"`
	BaseDenom          string    `toml:"base-denom"`
	DisplayDenom       string    `toml:"display-denom"`
	DenomCoefficient   int64     `toml:"denom-coefficient" default:"1000000"`
	LogUnknownMessages bool      `toml:"log-unknown-messages" default:"true"`

	Filters Filters
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

	for index, filter := range c.FiltersRaw {
		if err := ValidateFilter(filter); err != nil {
			return fmt.Errorf("Error in filter %d: %s", index, err)
		}
	}

	return nil
}

func ValidateFilter(filter string) error {
	split := strings.Split(filter, " ")
	if len(split) != 3 {
		return fmt.Errorf("filter should match pattern: <key> <operator> <value>")
	}

	if !utils.Contains([]string{"=", "!="}, split[1]) {
		return fmt.Errorf("unknown operator %s, allowed are: '=', '!='", split[1])
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
