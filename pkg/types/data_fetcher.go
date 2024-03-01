package types

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types/amount"
	"main/pkg/types/responses"

	transferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

// This interface is only here to avoid a cyclic dependency
// DataFetcher -> MetricsManager -> types -> DataFetcher.

type DataFetcher interface {
	GetPriceFetcher(info *configTypes.DenomInfo) PriceFetcher
	PopulateAmount(chainID string, amount *amount.Amount)
	PopulateAmounts(chainID string, amount amount.Amounts)
	GetRewardsAtBlock(
		chain *configTypes.Chain,
		delegator string,
		validator string,
		block int64,
	) ([]responses.Reward, bool)
	GetCommissionAtBlock(
		chain *configTypes.Chain,
		validator string,
		block int64,
	) ([]responses.Commission, bool)
	GetProposal(chain *configTypes.Chain, id string) (*responses.Proposal, bool)
	GetStakingParams(chain *configTypes.Chain) (*responses.StakingParams, bool)
	GetIbcRemoteChainID(chainID string, channel, port string) (string, bool)
	FindChainById(chainID string) (*configTypes.Chain, bool)
	GetDenomTrace(
		chain *configTypes.Chain,
		denom string,
	) (*transferTypes.DenomTrace, bool)
	PopulateWallet(chain *configTypes.Chain, walletLink *configTypes.Link, subscriptionName string)
	PopulateMultichainWallet(
		chain *configTypes.Chain,
		channel string,
		port string,
		walletLink *configTypes.Link,
		subscriptionName string,
	)
	PopulateWalletAlias(
		chain *configTypes.Chain,
		link *configTypes.Link,
		subscriptionName string,
	)
	PopulateValidator(
		chain *configTypes.Chain,
		validatorLink *configTypes.Link,
	)
	FindChainsByReporter(
		reporterName string,
	) configTypes.Chains
}
