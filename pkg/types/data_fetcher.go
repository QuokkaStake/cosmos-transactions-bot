package types

import (
	"main/pkg/alias_manager"
	configTypes "main/pkg/config/types"
	priceFetchers "main/pkg/price_fetchers"
	"main/pkg/types/amount"
	"main/pkg/types/responses"
)

// This interface is only here to avoid a cyclic dependencyL
// DataFetcher -> MetricsManager -> types -> DataFetcher.
type DataFetcher interface {
	GetPriceFetcher(info *configTypes.DenomInfo) priceFetchers.PriceFetcher
	PopulateAmount(amount *amount.Amount)
	PopulateAmounts(amount amount.Amounts)

	GetValidator(address string) (*responses.Validator, bool)
	GetRewardsAtBlock(
		delegator string,
		validator string,
		block int64,
	) ([]responses.Reward, bool)
	GetCommissionAtBlock(
		validator string,
		block int64,
	) ([]responses.Commission, bool)
	GetProposal(id string) (*responses.Proposal, bool)
	GetStakingParams() (*responses.StakingParams, bool)
	GetAliasManager() *alias_manager.AliasManager
	GetChain() *configTypes.Chain
}
