package alias_manager

import (
	configTypes "main/pkg/config/types"

	"github.com/rs/zerolog"
)

// Types that are used internally

// map[wallet]wallet_alias

type TomlChainAliases map[string]string

// map[chain]chain_aliases

type TomlSubscriptionAliases map[string]*TomlChainAliases

// map[subscription]subscription_aliases

type TomlAliases map[string]*TomlSubscriptionAliases

func (t TomlAliases) ToAliases(
	chains configTypes.Chains,
	logger zerolog.Logger,
) AllAliases {
	aliases := AllAliases{}

	for subscription, tomlSubscriptionAliases := range t {
		newMap := map[string]ChainAliases{}
		subscriptionAliases := &newMap

		for chain, tomlChainAliases := range *tomlSubscriptionAliases {
			chainFound := chains.FindByName(chain)
			if chainFound == nil {
				logger.Panic().Str("chain", chain).Msg("Could not find chain when setting an alias!")
			}

			chainAliases := make(map[string]string)

			for wallet, alias := range *tomlChainAliases {
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
