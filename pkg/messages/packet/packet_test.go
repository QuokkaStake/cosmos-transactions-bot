package packet

import (
	configTypes "main/pkg/config/types"
	"testing"

	types2 "github.com/cosmos/cosmos-sdk/codec/types"
	icaTypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcTypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcChannelTypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestPacketFungibleParseOk(t *testing.T) {
	t.Parallel()

	msg := &ibcTypes.FungibleTokenPacketData{
		Denom:    "uatom",
		Amount:   "100",
		Sender:   "sender",
		Receiver: "receiver",
		Memo:     "memo",
	}

	msgBytes, err := ibcTypes.ModuleCdc.MarshalJSON(msg)
	require.NoError(t, err)

	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "dst_channel",
		DestinationPort:    "dst_port",
		Data:               msgBytes,
	}

	parsed, err := ParsePacket(packet, &configTypes.Chain{Name: "chain"})
	require.NoError(t, err)
	require.NotNil(t, parsed)

	_, ok := parsed.(*FungibleTokenPacket)
	require.True(t, ok)
}

func TestPacketInterchainParseOk(t *testing.T) {
	t.Parallel()

	txInternal := &icaTypes.CosmosTx{
		Messages: []*types2.Any{},
	}
	txBytes, err := proto.Marshal(txInternal)
	require.NoError(t, err)

	icaPacket := &icaTypes.InterchainAccountPacketData{
		Type: icaTypes.EXECUTE_TX,
		Data: txBytes,
	}
	icaPacketBytes, err := ibcTypes.ModuleCdc.MarshalJSON(icaPacket)
	require.NoError(t, err)

	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "dst_channel",
		DestinationPort:    "dst_port",
		Data:               icaPacketBytes,
	}

	parsed, err := ParsePacket(packet, &configTypes.Chain{Name: "chain"})
	require.NoError(t, err)
	require.NotNil(t, parsed)

	_, ok := parsed.(*InterchainAccountsPacket)
	require.True(t, ok)
}

func TestPacketParseFail(t *testing.T) {
	t.Parallel()

	packet := ibcChannelTypes.Packet{
		SourceChannel:      "src_channel",
		SourcePort:         "src_port",
		DestinationChannel: "dst_channel",
		DestinationPort:    "dst_port",
		Data:               []byte("invalid"),
	}

	parsed, err := ParsePacket(packet, &configTypes.Chain{Name: "chain"})
	require.Error(t, err)
	require.Nil(t, parsed)
}
