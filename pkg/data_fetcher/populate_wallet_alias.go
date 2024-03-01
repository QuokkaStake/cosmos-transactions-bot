package data_fetcher

import configTypes "main/pkg/config/types"

func (f *DataFetcher) PopulateWalletAlias(
	chain *configTypes.Chain,
	link *configTypes.Link,
	subscriptionName string,
) {
	if alias := f.AliasManager.Get(subscriptionName, chain.Name, link.Value); alias != "" {
		link.Title = alias
	}
}
