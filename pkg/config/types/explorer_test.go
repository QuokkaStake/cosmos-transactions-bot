package types_test

import (
	"main/pkg/config/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMintscanExplorerToExplorer(t *testing.T) {
	t.Parallel()

	explorer := types.MintscanExplorer{Prefix: "chain"}
	appExplorer := explorer.ToExplorer()

	require.Equal(t, "https://mintscan.io/chain/blocks/%s", appExplorer.BlockLinkPattern)
	require.Equal(t, "https://mintscan.io/chain/tx/%s", appExplorer.TransactionLinkPattern)
	require.Equal(t, "https://mintscan.io/chain/validators/%s", appExplorer.ValidatorLinkPattern)
	require.Equal(t, "https://mintscan.io/chain/account/%s", appExplorer.WalletLinkPattern)
	require.Equal(t, "https://mintscan.io/chain/proposals/%s", appExplorer.ProposalLinkPattern)
}

func TestPingExplorerToExplorer(t *testing.T) {
	t.Parallel()

	explorer := types.PingExplorer{Prefix: "chain", BaseUrl: "https://example.com"}
	appExplorer := explorer.ToExplorer()

	require.Equal(t, "https://example.com/chain/blocks/%s", appExplorer.BlockLinkPattern)
	require.Equal(t, "https://example.com/chain/tx/%s", appExplorer.TransactionLinkPattern)
	require.Equal(t, "https://example.com/chain/staking/%s", appExplorer.ValidatorLinkPattern)
	require.Equal(t, "https://example.com/chain/account/%s", appExplorer.WalletLinkPattern)
	require.Equal(t, "https://example.com/chain/gov/%s", appExplorer.ProposalLinkPattern)
}
