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

	yamlAliases := alias_manager.YamlAliases{
		"subscription": &alias_manager.YamlSubscriptionAliases{
			"chain": &alias_manager.YamlChainAliases{
				"wallet": "alias",
			},
		},
	}

	logger := loggerPkg.GetNopLogger()
	chains := configTypes.Chains{
		{Name: "chain"},
	}

	aliases := yamlAliases.ToAliases(chains, *logger)
	require.Len(t, aliases, 1)

	subscriptionAliases := aliases["subscription"]
	require.Len(t, *subscriptionAliases, 1)

	chainAliases := (*subscriptionAliases)["chain"]
	require.Len(t, chainAliases.Aliases, 1)
	require.Equal(t, "alias", chainAliases.Aliases["wallet"])
}

func TestToAliasesNoChain(t *testing.T) {
	t.Parallel()

	yamlAliases := alias_manager.YamlAliases{
		"subscription": &alias_manager.YamlSubscriptionAliases{
			"chain": &alias_manager.YamlChainAliases{
				"wallet": "alias",
			},
		},
	}

	logger := loggerPkg.GetNopLogger()
	chains := configTypes.Chains{}

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	_ = yamlAliases.ToAliases(chains, *logger)
}
