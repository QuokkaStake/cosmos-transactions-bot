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

func TestExplorerDisplayWarningsNoTransactionLink(t *testing.T) {
	t.Parallel()

	explorer := &types.Explorer{
		BlockLinkPattern:     "test/%s",
		WalletLinkPattern:    "test/%s",
		ValidatorLinkPattern: "test/%s",
		ProposalLinkPattern:  "test/%s",
	}

	warnings := explorer.DisplayWarnings(&types.Chain{Name: "chain"})

	require.Len(t, warnings, 1)
}

func TestExplorerDisplayWarningsNoBlocksLink(t *testing.T) {
	t.Parallel()

	explorer := &types.Explorer{
		TransactionLinkPattern: "test/%s",
		WalletLinkPattern:      "test/%s",
		ValidatorLinkPattern:   "test/%s",
		ProposalLinkPattern:    "test/%s",
	}

	warnings := explorer.DisplayWarnings(&types.Chain{Name: "chain"})

	require.Len(t, warnings, 1)
}

func TestExplorerDisplayWarningsNoValidatorsLink(t *testing.T) {
	t.Parallel()

	explorer := &types.Explorer{
		TransactionLinkPattern: "test/%s",
		WalletLinkPattern:      "test/%s",
		BlockLinkPattern:       "test/%s",
		ProposalLinkPattern:    "test/%s",
	}

	warnings := explorer.DisplayWarnings(&types.Chain{Name: "chain"})

	require.Len(t, warnings, 1)
}

func TestExplorerDisplayWarningsNoWalletLink(t *testing.T) {
	t.Parallel()

	explorer := &types.Explorer{
		TransactionLinkPattern: "test/%s",
		ValidatorLinkPattern:   "test/%s",
		BlockLinkPattern:       "test/%s",
		ProposalLinkPattern:    "test/%s",
	}

	warnings := explorer.DisplayWarnings(&types.Chain{Name: "chain"})

	require.Len(t, warnings, 1)
}

func TestExplorerDisplayWarningsNoProposalLink(t *testing.T) {
	t.Parallel()

	explorer := &types.Explorer{
		TransactionLinkPattern: "test/%s",
		ValidatorLinkPattern:   "test/%s",
		BlockLinkPattern:       "test/%s",
		WalletLinkPattern:      "test/%s",
	}

	warnings := explorer.DisplayWarnings(&types.Chain{Name: "chain"})

	require.Len(t, warnings, 1)
}

func TestExplorerDisplayWarningsValid(t *testing.T) {
	t.Parallel()

	explorer := &types.Explorer{
		TransactionLinkPattern: "test/%s",
		BlockLinkPattern:       "test/%s",
		WalletLinkPattern:      "test/%s",
		ValidatorLinkPattern:   "test/%s",
		ProposalLinkPattern:    "test/%s",
	}

	warnings := explorer.DisplayWarnings(&types.Chain{Name: "chain"})

	require.Empty(t, warnings)
}
