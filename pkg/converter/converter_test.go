package converter_test

import (
	configTypes "main/pkg/config/types"
	converterPkg "main/pkg/converter"
	loggerPkg "main/pkg/logger"
	"main/pkg/messages"
	"main/pkg/types"
	"testing"

	jsonRpcTypes "github.com/cometbft/cometbft/rpc/jsonrpc/types"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	cosmosAuthzTypes "github.com/cosmos/cosmos-sdk/x/authz"
	cosmosBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

func TestConverterTxError(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	chain := &configTypes.Chain{Name: "chain"}
	converter := converterPkg.NewConverter(logger, chain)

	event := jsonRpcTypes.RPCResponse{
		Error: &jsonRpcTypes.RPCError{Message: "test message"},
	}
	result := converter.ParseEvent(event, "example")
	require.NotNil(t, result)
	require.IsType(t, &types.TxError{}, result)

	event2 := jsonRpcTypes.RPCResponse{
		Error: &jsonRpcTypes.RPCError{Message: "client is already subscribed"},
	}
	result2 := converter.ParseEvent(event2, "example")
	require.Nil(t, result2)
}

func TestConverterParseError(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	chain := &configTypes.Chain{Name: "chain"}
	converter := converterPkg.NewConverter(logger, chain)

	event := jsonRpcTypes.RPCResponse{
		Result: []byte("test"),
	}
	result := converter.ParseEvent(event, "example")
	require.NotNil(t, result)
	require.IsType(t, &types.TxError{}, result)
}

func TestConverterEmptyEvent(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	chain := &configTypes.Chain{Name: "chain"}
	converter := converterPkg.NewConverter(logger, chain)

	event := jsonRpcTypes.RPCResponse{
		Result: []byte("{\"data\":null}"),
	}
	result := converter.ParseEvent(event, "example")
	require.Nil(t, result)
}

func TestConverterErrorConvert(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	chain := &configTypes.Chain{Name: "chain"}
	converter := converterPkg.NewConverter(logger, chain)

	event := jsonRpcTypes.RPCResponse{
		Result: []byte("{\"data\":{\"type\":\"tendermint/event/NewBlock\",\"value\":{}}}"),
	}
	result := converter.ParseEvent(event, "example")
	require.Nil(t, result)
}

func TestConverterErrorUnmarshal(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	chain := &configTypes.Chain{Name: "chain"}
	converter := converterPkg.NewConverter(logger, chain)

	event := jsonRpcTypes.RPCResponse{
		Result: []byte("{\"data\":{\"type\":\"tendermint/event/Tx\",\"value\":{\"TxResult\":{\"height\":\"1\",\"index\":9,\"tx\":\"CmEKXwooL3NlbnRpbmVsLm5vZGUudjIuTXNnVXBkYXRlU3RhdHVzUmVxdWVzdBIzCi9zZW50bm9kZTFmdGNycnU0MDdmbGdhZTB0cm4wbjRja2RtY2w1aDZsZ3V4ZHIydBABEmcKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQIgzmHYcht/wBxkUOilsMRa2qUxPhHn8smOT1eFQBxlmRIECgIIARheEhMKDQoFdWR2cG4SBDk1MjYQk+gFGkB20FDj4l1Btj7avEltQAB3KH63PHg+52nXfcshadIwZmDErlv5dzF1Jz/d2NIs4gRj/5/twPFCabAffMlLsYlm\",\"result\":{\"data\":\"CisKKS9zZW50aW5lbC5ub2RlLnYyLk1zZ1VwZGF0ZURldGFpbHNSZXF1ZXN0\",\"log\":\"\",\"gas_wanted\":\"106365\",\"gas_used\":\"102726\",\"events\":[]}}}},\"events\":{}}"),
	}
	result := converter.ParseEvent(event, "example")
	require.NotNil(t, result)
	require.IsType(t, &types.Tx{}, result)
}

func TestConverterUnsupportedMessage(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	chain := &configTypes.Chain{Name: "chain"}
	converter := converterPkg.NewConverter(logger, chain)

	message := &codecTypes.Any{TypeUrl: "unsupported"}
	result := converter.ParseMessage(message, 123)
	require.NotNil(t, result)
	require.IsType(t, &messages.MsgUnsupportedMessage{}, result)
}

func TestConverterUnparsedMessage(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	chain := &configTypes.Chain{Name: "chain"}
	converter := converterPkg.NewConverter(logger, chain)

	message := &codecTypes.Any{
		TypeUrl: "/cosmos.bank.v1beta1.MsgSend",
		Value:   []byte("unparsed"),
	}
	result := converter.ParseMessage(message, 123)
	require.NotNil(t, result)
	require.IsType(t, &messages.MsgUnparsedMessage{}, result)
}

func TestConverterParsedCorrectly(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	chain := &configTypes.Chain{Name: "chain"}
	converter := converterPkg.NewConverter(logger, chain)

	msgSend := &cosmosBankTypes.MsgSend{}
	bytes, err := msgSend.Marshal()
	require.NoError(t, err)

	message := &codecTypes.Any{
		TypeUrl: "/cosmos.bank.v1beta1.MsgSend",
		Value:   bytes,
	}
	result := converter.ParseMessage(message, 123)
	require.NotNil(t, result)
	require.IsType(t, &messages.MsgSend{}, result)
}

func TestConverterParsedInternal(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	chain := &configTypes.Chain{Name: "chain"}
	converter := converterPkg.NewConverter(logger, chain)

	msgSend := &cosmosBankTypes.MsgSend{}
	msgSendBytes, err := msgSend.Marshal()
	require.NoError(t, err)

	msgExec := &cosmosAuthzTypes.MsgExec{
		Msgs: []*codecTypes.Any{
			{
				TypeUrl: "/cosmos.bank.v1beta1.MsgSend",
				Value:   msgSendBytes,
			},
		},
	}
	bytes, err := msgExec.Marshal()
	require.NoError(t, err)

	message := &codecTypes.Any{
		TypeUrl: "/cosmos.authz.v1beta1.MsgExec",
		Value:   bytes,
	}
	result := converter.ParseMessage(message, 123)
	require.NotNil(t, result)
	require.IsType(t, &messages.MsgExec{}, result)
	require.Len(t, result.GetParsedMessages(), 1)
	require.IsType(t, &messages.MsgSend{}, result.GetParsedMessages()[0])
}
