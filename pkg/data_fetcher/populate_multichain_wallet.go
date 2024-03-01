package data_fetcher

import (
	configTypes "main/pkg/config/types"
)

func (f *DataFetcher) PopulateWallet(
	chain *configTypes.Chain,
	walletLink *configTypes.Link,
	subscriptionName string,
) {
	if chain.Explorer == nil {
		return
	}

	walletLink.Href = chain.Explorer.GetWalletLink(walletLink.Value)
	if alias := f.AliasManager.Get(subscriptionName, chain.Name, walletLink.Value); alias != "" {
		walletLink.Title = alias
	}
}

func (f *DataFetcher) PopulateMultichainWallet(
	chain *configTypes.Chain,
	channel string,
	port string,
	walletLink *configTypes.Link,
	subscriptionName string,
) {
	// Wallet from local chain, take local explorer config.
	if channel == "" || port == "" {
		f.PopulateWallet(chain, walletLink, subscriptionName)
		return
	}

	// Wallet is from another chain. Resolving its chain-id it by traversing the IBC path.
	remoteChainId, fetched := f.GetIbcRemoteChainID(chain.ChainID, channel, port)
	if !fetched {
		return
	}

	// Trying to find it in our config.
	localChain, found := f.FindChainById(remoteChainId)
	if !found {
		return
	}

	if localChain.Explorer == nil {
		return
	}

	walletLink.Href = localChain.Explorer.GetWalletLink(walletLink.Value)
	if alias := f.AliasManager.Get(subscriptionName, chain.Name, walletLink.Value); alias != "" {
		walletLink.Title = alias
	}
}
