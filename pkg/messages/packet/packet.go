package packet

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types"

	icaTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcChannelTypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

func ParsePacket(message ibcChannelTypes.Packet, chain *configTypes.Chain) (types.Message, error) {
	// Fungible token transfer
	var fungiblePacketData ibcTypes.FungibleTokenPacketData
	err := ibcTypes.ModuleCdc.UnmarshalJSON(message.Data, &fungiblePacketData)
	if err == nil {
		return ParseFungibleTokenPacket(fungiblePacketData, message, chain), nil
	}

	// ICA packet
	var icaPacketData icaTypes.InterchainAccountPacketData
	err = ibcTypes.ModuleCdc.UnmarshalJSON(message.Data, &icaPacketData)
	if err == nil {
		return ParseInterchainAccountsPacket(icaPacketData, chain)
	}

	return nil, err
}
