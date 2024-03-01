package alias_manager

import (
	configTypes "main/pkg/config/types"
)

// Types that are used to save/load TOML.

type ChainAliases struct {
	Chain *configTypes.Chain

	// map[wallet]wallet_alias
	Aliases map[string]string
}

// map[chain]chain_aliases

type SubscriptionAliases *map[string]ChainAliases

// map[subscription]subscription_aliases

type AllAliases map[string]SubscriptionAliases

func (a AllAliases) ToTomlAliases() *TomlAliases {
	tomlAliases := TomlAliases{}

	for subscription, subscriptionAliases := range a {
		tomlSubscriptionAliases := TomlSubscriptionAliases{}

		for chain, chainAliases := range *subscriptionAliases {
			tomlChainAliases := TomlChainAliases{}
			for wallet, alias := range chainAliases.Aliases {
				tomlChainAliases[wallet] = alias
			}

			tomlSubscriptionAliases[chain] = &tomlChainAliases
		}

		tomlAliases[subscription] = &tomlSubscriptionAliases
	}

	return &tomlAliases
}

func (a AllAliases) Set(subscription string, chain *configTypes.Chain, wallet, alias string) {
	if _, ok := a[subscription]; !ok {
		newMap := map[string]ChainAliases{}
		a[subscription] = &newMap
	}

	subscriptionAliases := a[subscription]

	if _, ok := (*subscriptionAliases)[chain.Name]; !ok {
		(*subscriptionAliases)[chain.Name] = ChainAliases{
			Chain:   chain,
			Aliases: make(map[string]string),
		}
	}

	chainAliases := (*subscriptionAliases)[chain.Name]
	chainAliases.Aliases[wallet] = alias
}

func (a AllAliases) Get(subscription, chain, address string) string {
	subscriptionAliases, ok := a[subscription]
	if !ok {
		return ""
	}

	chainAliases, ok := (*subscriptionAliases)[chain]
	if !ok {
		return ""
	}

	aliases := chainAliases.Aliases
	alias, ok := aliases[address]
	if !ok {
		return ""
	}

	return alias
}

type ChainAliasesLinks struct {
	Chain *configTypes.Chain
	Links map[string]*configTypes.Link
}

func (a AllAliases) GetAliasesLinks(subscription string) []ChainAliasesLinks {
	subscriptionAliases, ok := a[subscription]
	if !ok {
		return []ChainAliasesLinks{}
	}

	aliasesLinks := make([]ChainAliasesLinks, 0)

	for _, chainAliases := range *subscriptionAliases {
		links := make(map[string]*configTypes.Link)

		for wallet, alias := range chainAliases.Aliases {
			link := chainAliases.Chain.GetWalletLink(wallet)
			link.Title = alias
			links[wallet] = link
		}

		aliasesLinks = append(aliasesLinks, ChainAliasesLinks{
			Chain: chainAliases.Chain,
			Links: links,
		})
	}

	return aliasesLinks
}
