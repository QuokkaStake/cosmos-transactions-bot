package yaml_config_test

import (
	"main/pkg/config/types"
	yamlConfig "main/pkg/config/yaml_config"
	"testing"

	queryPkg "github.com/cometbft/cometbft/libs/pubsub/query"
	"github.com/stretchr/testify/require"
)

func TestChainEmptyName(t *testing.T) {
	t.Parallel()

	chain := yamlConfig.Chain{}
	require.Error(t, chain.Validate())
}
func TestChainEmptyID(t *testing.T) {
	t.Parallel()

	chain := yamlConfig.Chain{
		Name: "chain",
	}
	require.Error(t, chain.Validate())
}

func TestChainEmptyTendermintNodes(t *testing.T) {
	t.Parallel()

	chain := yamlConfig.Chain{
		Name:    "chain",
		ChainID: "chain-id",
	}
	require.Error(t, chain.Validate())
}

func TestChainEmptyApiNodes(t *testing.T) {
	t.Parallel()

	chain := yamlConfig.Chain{
		Name:            "chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"node"},
	}
	require.Error(t, chain.Validate())
}

func TestChainEmptyQueries(t *testing.T) {
	t.Parallel()

	chain := yamlConfig.Chain{
		Name:            "chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
	}
	require.Error(t, chain.Validate())
}

func TestChainInvalidQuery(t *testing.T) {
	t.Parallel()

	chain := yamlConfig.Chain{
		Name:            "chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"query"},
	}
	require.Error(t, chain.Validate())
}

func TestChainInvalidDenom(t *testing.T) {
	t.Parallel()

	chain := yamlConfig.Chain{
		Name:            "chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
		Denoms: yamlConfig.DenomInfos{
			{},
		},
	}
	require.Error(t, chain.Validate())
}

func TestChainValid(t *testing.T) {
	t.Parallel()

	chain := yamlConfig.Chain{
		Name:            "chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
	}
	require.NoError(t, chain.Validate())
}

func TestChainToAppConfigChainBasic(t *testing.T) {
	t.Parallel()

	chain := yamlConfig.Chain{
		Name:            "chain",
		PrettyName:      "Chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"tendermint-node"},
		APINodes:        []string{"api-node"},
		Queries:         []string{"event.key = 'value'"},
	}
	appConfigChain := chain.ToAppConfigChain()

	require.Equal(t, "chain", appConfigChain.Name)
	require.Equal(t, "Chain", appConfigChain.PrettyName)
	require.Equal(t, "chain-id", appConfigChain.ChainID)
	require.Len(t, appConfigChain.TendermintNodes, 1)
	require.Equal(t, "tendermint-node", appConfigChain.TendermintNodes[0])
	require.Len(t, appConfigChain.APINodes, 1)
	require.Equal(t, "api-node", appConfigChain.APINodes[0])
	require.Len(t, appConfigChain.Queries, 1)
	require.Equal(t, "event.key = 'value'", appConfigChain.Queries[0].String())
}

func TestChainToAppConfigChainMintscan(t *testing.T) {
	t.Parallel()

	chain := yamlConfig.Chain{
		Name:            "chain",
		PrettyName:      "Chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"tendermint-node"},
		APINodes:        []string{"api-node"},
		Queries:         []string{"event.key = 'value'"},
		MintscanPrefix:  "chain",
	}
	appConfigChain := chain.ToAppConfigChain()

	require.NotNil(t, appConfigChain.Explorer)
	require.Equal(t, "https://mintscan.io/chain/validators/%s", appConfigChain.Explorer.ValidatorLinkPattern)
}

func TestChainToAppConfigChainPing(t *testing.T) {
	t.Parallel()

	chain := yamlConfig.Chain{
		Name:            "chain",
		PrettyName:      "Chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"tendermint-node"},
		APINodes:        []string{"api-node"},
		Queries:         []string{"event.key = 'value'"},
		PingPrefix:      "chain",
		PingBaseUrl:     "https://example.com",
	}
	appConfigChain := chain.ToAppConfigChain()

	require.NotNil(t, appConfigChain.Explorer)
	require.Equal(t, "https://example.com/chain/staking/%s", appConfigChain.Explorer.ValidatorLinkPattern)
}

func TestChainToAppConfigChainCustomExplorer(t *testing.T) {
	t.Parallel()

	chain := yamlConfig.Chain{
		Name:            "chain",
		PrettyName:      "Chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"tendermint-node"},
		APINodes:        []string{"api-node"},
		Queries:         []string{"event.key = 'value'"},
		Explorer: &yamlConfig.Explorer{
			ValidatorLinkPattern: "test/%s",
		},
	}
	appConfigChain := chain.ToAppConfigChain()

	require.NotNil(t, appConfigChain.Explorer)
	require.Equal(t, "test/%s", appConfigChain.Explorer.ValidatorLinkPattern)
}

func TestChainToYamlConfigChainBasic(t *testing.T) {
	t.Parallel()

	query := queryPkg.MustParse("event.key = 'value'")

	chain := &types.Chain{
		Name:            "chain",
		PrettyName:      "Chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"tendermint-node"},
		APINodes:        []string{"api-node"},
		Queries:         []queryPkg.Query{*query},
	}

	yamlConfigChain := yamlConfig.FromAppConfigChain(chain)

	require.Equal(t, "chain", yamlConfigChain.Name)
	require.Equal(t, "Chain", yamlConfigChain.PrettyName)
	require.Equal(t, "chain-id", yamlConfigChain.ChainID)
	require.Len(t, yamlConfigChain.TendermintNodes, 1)
	require.Equal(t, "tendermint-node", yamlConfigChain.TendermintNodes[0])
	require.Len(t, yamlConfigChain.APINodes, 1)
	require.Equal(t, "api-node", yamlConfigChain.APINodes[0])
	require.Len(t, yamlConfigChain.Queries, 1)
	require.Equal(t, "event.key = 'value'", yamlConfigChain.Queries[0])
}

func TestChainToYamlConfigChainMintscan(t *testing.T) {
	t.Parallel()

	query := queryPkg.MustParse("event.key = 'value'")

	chain := &types.Chain{
		Name:              "chain",
		PrettyName:        "Chain",
		ChainID:           "chain-id",
		TendermintNodes:   []string{"tendermint-node"},
		APINodes:          []string{"api-node"},
		Queries:           []queryPkg.Query{*query},
		SupportedExplorer: &types.MintscanExplorer{Prefix: "chain"},
	}
	yamlConfigChain := yamlConfig.FromAppConfigChain(chain)

	require.Equal(t, "chain", yamlConfigChain.MintscanPrefix)
}

func TestChainToYamlConfigChainPing(t *testing.T) {
	t.Parallel()

	query := queryPkg.MustParse("event.key = 'value'")

	chain := &types.Chain{
		Name:              "chain",
		PrettyName:        "Chain",
		ChainID:           "chain-id",
		TendermintNodes:   []string{"tendermint-node"},
		APINodes:          []string{"api-node"},
		Queries:           []queryPkg.Query{*query},
		SupportedExplorer: &types.PingExplorer{Prefix: "chain", BaseUrl: "https://example.com"},
	}
	yamlConfigChain := yamlConfig.FromAppConfigChain(chain)

	require.Equal(t, "chain", yamlConfigChain.PingPrefix)
	require.Equal(t, "https://example.com", yamlConfigChain.PingBaseUrl)
}

func TestChainToYamlConfigChainCustomExplorer(t *testing.T) {
	t.Parallel()

	query := queryPkg.MustParse("event.key = 'value'")

	chain := &types.Chain{
		Name:            "chain",
		PrettyName:      "Chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"tendermint-node"},
		APINodes:        []string{"api-node"},
		Queries:         []queryPkg.Query{*query},
		Explorer: &types.Explorer{
			ValidatorLinkPattern: "test/%s",
		},
	}
	yamlConfigChain := yamlConfig.FromAppConfigChain(chain)

	require.NotNil(t, yamlConfigChain.Explorer)
	require.Equal(t, "test/%s", yamlConfigChain.Explorer.ValidatorLinkPattern)
}

func TestChainsInvalid(t *testing.T) {
	t.Parallel()

	chain := &yamlConfig.Chain{}
	chains := yamlConfig.Chains{chain}

	require.Error(t, chains.Validate())
}

func TestChainsDuplicateName(t *testing.T) {
	t.Parallel()

	chain1 := &yamlConfig.Chain{
		Name:            "chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
	}
	chain2 := &yamlConfig.Chain{
		Name:            "chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
	}
	chains := yamlConfig.Chains{chain1, chain2}

	require.Error(t, chains.Validate())
}

func TestChainsValid(t *testing.T) {
	t.Parallel()

	chain1 := &yamlConfig.Chain{
		Name:            "chain1",
		ChainID:         "chain-id",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
	}
	chain2 := &yamlConfig.Chain{
		Name:            "chain2",
		ChainID:         "chain-id",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
	}
	chains := yamlConfig.Chains{chain1, chain2}

	require.NoError(t, chains.Validate())
}

func TestHasChainByName(t *testing.T) {
	t.Parallel()

	chain := &yamlConfig.Chain{
		Name:            "chain-1",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
	}
	chains := yamlConfig.Chains{chain}

	require.True(t, chains.HasChainByName("chain-1"))
	require.False(t, chains.HasChainByName("chain-2"))
}
