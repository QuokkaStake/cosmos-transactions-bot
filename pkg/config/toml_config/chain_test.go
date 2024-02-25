package toml_config_test

import (
	tomlConfig "main/pkg/config/toml_config"
	"main/pkg/config/types"
	"testing"

	queryPkg "github.com/cometbft/cometbft/libs/pubsub/query"
	"github.com/stretchr/testify/require"
)

func TestChainEmptyName(t *testing.T) {
	t.Parallel()

	chain := tomlConfig.Chain{}
	require.Error(t, chain.Validate())
}

func TestChainEmptyTendermintNodes(t *testing.T) {
	t.Parallel()

	chain := tomlConfig.Chain{
		Name: "chain",
	}
	require.Error(t, chain.Validate())
}

func TestChainEmptyApiNodes(t *testing.T) {
	t.Parallel()

	chain := tomlConfig.Chain{
		Name:            "chain",
		TendermintNodes: []string{"node"},
	}
	require.Error(t, chain.Validate())
}

func TestChainEmptyQueries(t *testing.T) {
	t.Parallel()

	chain := tomlConfig.Chain{
		Name:            "chain",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
	}
	require.Error(t, chain.Validate())
}

func TestChainInvalidQuery(t *testing.T) {
	t.Parallel()

	chain := tomlConfig.Chain{
		Name:            "chain",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"query"},
	}
	require.Error(t, chain.Validate())
}

func TestChainInvalidDenom(t *testing.T) {
	t.Parallel()

	chain := tomlConfig.Chain{
		Name:            "chain",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
		Denoms: tomlConfig.DenomInfos{
			{},
		},
	}
	require.Error(t, chain.Validate())
}

func TestChainValid(t *testing.T) {
	t.Parallel()

	chain := tomlConfig.Chain{
		Name:            "chain",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
	}
	require.NoError(t, chain.Validate())
}

func TestChainToAppConfigChainBasic(t *testing.T) {
	t.Parallel()

	chain := tomlConfig.Chain{
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

	chain := tomlConfig.Chain{
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

	chain := tomlConfig.Chain{
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

	chain := tomlConfig.Chain{
		Name:            "chain",
		PrettyName:      "Chain",
		ChainID:         "chain-id",
		TendermintNodes: []string{"tendermint-node"},
		APINodes:        []string{"api-node"},
		Queries:         []string{"event.key = 'value'"},
		Explorer: &tomlConfig.Explorer{
			ValidatorLinkPattern: "test/%s",
		},
	}
	appConfigChain := chain.ToAppConfigChain()

	require.NotNil(t, appConfigChain.Explorer)
	require.Equal(t, "test/%s", appConfigChain.Explorer.ValidatorLinkPattern)
}

func TestChainToTomlConfigChainBasic(t *testing.T) {
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

	tomlConfigChain := tomlConfig.FromAppConfigChain(chain)

	require.Equal(t, "chain", tomlConfigChain.Name)
	require.Equal(t, "Chain", tomlConfigChain.PrettyName)
	require.Equal(t, "chain-id", tomlConfigChain.ChainID)
	require.Len(t, tomlConfigChain.TendermintNodes, 1)
	require.Equal(t, "tendermint-node", tomlConfigChain.TendermintNodes[0])
	require.Len(t, tomlConfigChain.APINodes, 1)
	require.Equal(t, "api-node", tomlConfigChain.APINodes[0])
	require.Len(t, tomlConfigChain.Queries, 1)
	require.Equal(t, "event.key = 'value'", tomlConfigChain.Queries[0])
}

func TestChainToTomlConfigChainMintscan(t *testing.T) {
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
	tomlConfigChain := tomlConfig.FromAppConfigChain(chain)

	require.Equal(t, "chain", tomlConfigChain.MintscanPrefix)
}

func TestChainToTomlConfigChainPing(t *testing.T) {
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
	tomlConfigChain := tomlConfig.FromAppConfigChain(chain)

	require.Equal(t, "chain", tomlConfigChain.PingPrefix)
	require.Equal(t, "https://example.com", tomlConfigChain.PingBaseUrl)
}

func TestChainToTomlConfigChainCustomExplorer(t *testing.T) {
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
	tomlConfigChain := tomlConfig.FromAppConfigChain(chain)

	require.NotNil(t, tomlConfigChain.Explorer)
	require.Equal(t, "test/%s", tomlConfigChain.Explorer.ValidatorLinkPattern)
}

func TestChainsInvalid(t *testing.T) {
	t.Parallel()

	chain := &tomlConfig.Chain{}
	chains := tomlConfig.Chains{chain}

	require.Error(t, chains.Validate())
}

func TestChainsDuplicateName(t *testing.T) {
	t.Parallel()

	chain1 := &tomlConfig.Chain{
		Name:            "chain",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
	}
	chain2 := &tomlConfig.Chain{
		Name:            "chain",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
	}
	chains := tomlConfig.Chains{chain1, chain2}

	require.Error(t, chains.Validate())
}

func TestChainsValid(t *testing.T) {
	t.Parallel()

	chain1 := &tomlConfig.Chain{
		Name:            "chain1",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
	}
	chain2 := &tomlConfig.Chain{
		Name:            "chain2",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
	}
	chains := tomlConfig.Chains{chain1, chain2}

	require.NoError(t, chains.Validate())
}

func TestHasChainByName(t *testing.T) {
	t.Parallel()

	chain := &tomlConfig.Chain{
		Name:            "chain-1",
		TendermintNodes: []string{"node"},
		APINodes:        []string{"node"},
		Queries:         []string{"event.key = 'value'"},
	}
	chains := tomlConfig.Chains{chain}

	require.True(t, chains.HasChainByName("chain-1"))
	require.False(t, chains.HasChainByName("chain-2"))
}
