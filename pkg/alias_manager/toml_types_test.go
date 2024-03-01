package alias_manager_test

import (
	"main/pkg/alias_manager"
	configTypes "main/pkg/config/types"
	loggerPkg "main/pkg/logger"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToAliasesValid(t *testing.T) {
	t.Parallel()

	tomlAliases := alias_manager.TomlAliases{
		"subscription": &alias_manager.TomlSubscriptionAliases{
			"chain": &alias_manager.TomlChainAliases{
				"wallet": "alias",
			},
		},
	}

	logger := loggerPkg.GetDefaultLogger()
	chains := configTypes.Chains{
		{Name: "chain"},
	}

	aliases := tomlAliases.ToAliases(chains, *logger)
	require.Len(t, aliases, 1)

	subscriptionAliases := aliases["subscription"]
	require.Len(t, *subscriptionAliases, 1)

	chainAliases := (*subscriptionAliases)["chain"]
	require.Len(t, chainAliases.Aliases, 1)
	require.Equal(t, "alias", chainAliases.Aliases["wallet"])
}

func TestToAliasesNoChain(t *testing.T) {
	t.Parallel()

	tomlAliases := alias_manager.TomlAliases{
		"subscription": &alias_manager.TomlSubscriptionAliases{
			"chain": &alias_manager.TomlChainAliases{
				"wallet": "alias",
			},
		},
	}

	logger := loggerPkg.GetDefaultLogger()
	chains := configTypes.Chains{}

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	_ = tomlAliases.ToAliases(chains, *logger)
}
