package data_fetcher

import (
	configTypes "main/pkg/config/types"
	QueryInfo "main/pkg/types/query_info"
	"main/pkg/types/responses"
)

func (f *DataFetcher) GetIbcRemoteChainID(
	chain *configTypes.Chain,
	channel string,
	port string,
) (string, bool) {
	keyName := chain.Name + "_channel_" + channel + "_port_" + port

	if cachedEntry, cachedEntryPresent := f.Cache.Get(keyName); cachedEntryPresent {
		if cachedEntryParsed, ok := cachedEntry.(string); ok {
			return cachedEntryParsed, true
		}

		f.Logger.Error().Msg("Could not convert cached IBC channel to string")
		return "", false
	}

	var (
		ibcChannel     *responses.IbcChannel
		ibcClientState *responses.IbcIdentifiedClientState
	)

	for _, node := range f.TendermintApiClients[chain.Name] {
		ibcChannelResponse, err, queryInfo := node.GetIbcChannel(channel, port)
		f.MetricsManager.LogTendermintQuery(chain.Name, queryInfo, QueryInfo.QueryTypeIbcChannel)
		if err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching IBC channel")
			continue
		}

		ibcChannel = ibcChannelResponse
		break
	}

	if ibcChannel == nil {
		f.Logger.Error().Msg("Could not connect to any nodes to get IBC channel")
		return "", false
	}

	if len(ibcChannel.ConnectionHops) != 1 {
		f.Logger.Error().
			Int("len", len(ibcChannel.ConnectionHops)).
			Msg("Could not connect to any nodes to get IBC channel")
		return "", false
	}

	for _, node := range f.TendermintApiClients[chain.Name] {
		ibcChannelClientStateResponse, err, queryInfo := node.GetIbcConnectionClientState(ibcChannel.ConnectionHops[0])
		f.MetricsManager.LogTendermintQuery(chain.Name, queryInfo, QueryInfo.QueryTypeIbcConnectionClientState)
		if err != nil {
			f.Logger.Error().Err(err).Msg("Error fetching IBC client state")
			continue
		}

		ibcClientState = ibcChannelClientStateResponse
		break
	}

	if ibcClientState == nil {
		f.Logger.Error().Msg("Could not connect to any nodes to get IBC client state")
		return "", false
	}

	f.Cache.Set(keyName, ibcClientState.ClientState.ChainId)
	return ibcClientState.ClientState.ChainId, true
}
