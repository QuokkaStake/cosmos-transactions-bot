package alias_manager

import (
	configTypes "main/pkg/config/types"

	"github.com/rs/zerolog"
)

// Types that are used internally

// map[wallet]wallet_alias

type YamlChainAliases map[string]string

// map[chain]chain_aliases

type YamlSubscriptionAliases map[string]*YamlChainAliases

// map[subscription]subscription_aliases

type YamlAliases map[string]*YamlSubscriptionAliases

func (t YamlAliases) ToAliases(
	chains configTypes.Chains,
	logger zerolog.Logger,
) AllAliases {
	aliases := AllAliases{}

	for subscription, yamlSubscriptionAliases := range t {
		newMap := map[string]ChainAliases{}
		subscriptionAliases := &newMap

		for chain, yamlChainAliases := range *yamlSubscriptionAliases {
			chainFound := chains.FindByName(chain)
			if chainFound == nil {
				logger.Panic().Str("chain", chain).Msg("Could not find chain when setting an alias!")
			}

			chainAliases := make(map[string]string)

			for wallet, alias := range *yamlChainAliases {
				chainAliases[wallet] = alias
			}

			(*subscriptionAliases)[chain] = ChainAliases{
				Chain:   chainFound,
				Aliases: chainAliases,
			}
		}

		aliases[subscription] = subscriptionAliases
	}

	return aliases
}
