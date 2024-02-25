package toml_config_test

import (
	tomlConfig "main/pkg/config/toml_config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExplorerToAppConfigExplorer(t *testing.T) {
	t.Parallel()

	explorer := &tomlConfig.Explorer{
		ValidatorLinkPattern:   "test1",
		WalletLinkPattern:      "test2",
		ProposalLinkPattern:    "test3",
		TransactionLinkPattern: "test4",
		BlockLinkPattern:       "test5",
	}
	appConfigExplorer := explorer.ToAppConfigExplorer()

	require.Equal(t, "test1", appConfigExplorer.ValidatorLinkPattern)
	require.Equal(t, "test2", appConfigExplorer.WalletLinkPattern)
	require.Equal(t, "test3", appConfigExplorer.ProposalLinkPattern)
	require.Equal(t, "test4", appConfigExplorer.TransactionLinkPattern)
	require.Equal(t, "test5", appConfigExplorer.BlockLinkPattern)
}
