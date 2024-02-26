package data_fetcher

import (
	configTypes "main/pkg/config/types"
	"strings"

	transferTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

func (f *DataFetcher) GetDenomTrace(
	chain *configTypes.Chain,
	denom string,
) (*transferTypes.DenomTrace, bool) {
	denomSplit := strings.Split(denom, "/")
	if len(denomSplit) != 2 || denomSplit[0] != transferTypes.DenomPrefix {
		f.Logger.Error().Msg("Invalid IBC prefix provided")
		return nil, false
	}

	denomHash := denomSplit[1]

	keyName := chain.Name + "_denom_trace_" + denom

	if cachedEntry, cachedEntryPresent := f.Cache.Get(keyName); cachedEntryPresent {
		if cachedEntryParsed, ok := cachedEntry.(*transferTypes.DenomTrace); ok {
			return cachedEntryParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached staking params to *types.DenomTrace")
		return nil, false
	}

	for _, node := range f.TendermintApiClients[chain.Name] {
		notCachedEntry, err := node.GetIbcDenomTrace(denomHash)

		if err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching IBC denom trace")
			continue
		}

		f.Cache.Set(keyName, notCachedEntry)
		return notCachedEntry, true
	}

	f.Logger.Error().Msg("Could not connect to any nodes to get IBC denom trace")
	return nil, false
}
