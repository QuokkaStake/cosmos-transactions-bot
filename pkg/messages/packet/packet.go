package packet

import (
	icaTypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	ibcTypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibcChannelTypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	configTypes "main/pkg/config/types"
	"main/pkg/types"
)

func ParsePacket(packet ibcChannelTypes.Packet, chain *configTypes.Chain) (types.Message, error) {
	// Fungible token transfer
	var fungiblePacketData ibcTypes.FungibleTokenPacketData
	err := ibcTypes.ModuleCdc.UnmarshalJSON(packet.Data, &fungiblePacketData)
	if err == nil {
		return ParseFungibleTokenPacket(fungiblePacketData, chain), err
	}

	// ICA packet
	var icaPacketData icaTypes.InterchainAccountPacketData
	err = ibcTypes.ModuleCdc.UnmarshalJSON(packet.Data, &icaPacketData)
	if err == nil {
		return ParseInterchainAccountsPacket(icaPacketData, chain)
	}

	return nil, err
}
