package converter

import (
	"fmt"
	configTypes "main/pkg/config/types"
	"main/pkg/messages"
	"main/pkg/types"
	"strings"

	"github.com/cometbft/cometbft/crypto/tmhash"
	"github.com/cometbft/cometbft/libs/json"
	coreTypes "github.com/cometbft/cometbft/rpc/core/types"
	jsonRpcTypes "github.com/cometbft/cometbft/rpc/jsonrpc/types"
	tendermintTypes "github.com/cometbft/cometbft/types"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog"
)

type Converter struct {
	Logger  zerolog.Logger
	Chain   *configTypes.Chain
	Parsers map[string]types.MessageParser
}

func NewConverter(logger *zerolog.Logger, chain *configTypes.Chain) *Converter {
	parsers := map[string]types.MessageParser{
		"/cosmos.authz.v1beta1.MsgExec":                               messages.ParseMsgExec,
		"/cosmos.authz.v1beta1.MsgGrant":                              messages.ParseMsgGrant,
		"/cosmos.authz.v1beta1.MsgRevoke":                             messages.ParseMsgRevoke,
		"/cosmos.bank.v1beta1.MsgSend":                                messages.ParseMsgSend,
		"/cosmos.bank.v1beta1.MsgMultiSend":                           messages.ParseMsgMultiSend,
		"/cosmos.distribution.v1beta1.MsgSetWithdrawAddress":          messages.ParseMsgSetWithdrawAddress,
		"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward":     messages.ParseMsgWithdrawDelegatorReward,
		"/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission": messages.ParseMsgWithdrawValidatorCommission,
		"/cosmos.gov.v1beta1.MsgVote":                                 messages.ParseMsgVote,
		"/cosmos.staking.v1beta1.MsgDelegate":                         messages.ParseMsgDelegate,
		"/cosmos.staking.v1beta1.MsgBeginRedelegate":                  messages.ParseMsgBeginRedelegate,
		"/cosmos.staking.v1beta1.MsgUndelegate":                       messages.ParseMsgUndelegate,
		"/ibc.applications.transfer.v1.MsgTransfer":                   messages.ParseMsgTransfer,
		"/ibc.core.channel.v1.MsgAcknowledgement":                     messages.ParseMsgAcknowledgement,
		"/ibc.core.channel.v1.MsgRecvPacket":                          messages.ParseMsgRecvPacket,
		"/ibc.core.channel.v1.MsgTimeout":                             messages.ParseMsgTimeout,
		"/ibc.core.client.v1.MsgUpdateClient":                         messages.ParseMsgUpdateClient,
	}

	return &Converter{
		Logger: logger.With().
			Str("component", "converter").
			Str("chain", chain.Name).
			Logger(),
		Parsers: parsers,
		Chain:   chain,
	}
}

func (c *Converter) ParseEvent(event jsonRpcTypes.RPCResponse, nodeURL string) types.Reportable {
	if event.Error != nil && event.Error.Message != "" {
		if strings.Contains(event.Error.Error(), "already subscribed") {
			c.Logger.Error().
				Str("msg", event.Error.Error()).
				Str("node", nodeURL).
				Msg("Client is already subscribed!")
			return nil
		}

		c.Logger.Error().
			Str("msg", event.Error.Error()).
			Str("node", nodeURL).
			Msg("Got error in RPC response")
		return &types.TxError{Error: event.Error}
	}

	var resultEvent coreTypes.ResultEvent
	if err := json.Unmarshal(event.Result, &resultEvent); err != nil {
		c.Logger.Error().Err(err).Msg("Failed to parse event")
		return &types.TxError{Error: event.Error}
	}

	if resultEvent.Data == nil {
		c.Logger.Debug().Msg("Event does not have data, skipping.")
		return nil
	}

	c.Logger.Trace().
		Str("values", fmt.Sprintf("%+v", resultEvent.Events)).
		Msg("Event values")

	eventDataTx, ok := resultEvent.Data.(tendermintTypes.EventDataTx)
	if !ok {
		c.Logger.Debug().Msg("Could not convert tx result to EventDataTx.")
		return nil
	}

	txResult := eventDataTx.TxResult
	txHash := fmt.Sprintf("%X", tmhash.Sum(txResult.Tx))

	var txProto tx.Tx

	if err := proto.Unmarshal(txResult.Tx, &txProto); err != nil {
		c.Logger.Error().Err(err).Msg("Could not parse tx")
		return &types.TxError{Error: event.Error}
	}

	c.Logger.Debug().
		Int64("height", txResult.Height).
		Str("memo", txProto.GetBody().GetMemo()).
		Str("hash", txHash).
		Int("len", len(txProto.GetBody().Messages)).
		Str("node", nodeURL).
		Msg("Got transaction")

	txMessages := []types.Message{}

	for _, message := range txProto.GetBody().Messages {
		if msgParsed := c.ParseMessage(message, txResult.Height); msgParsed != nil {
			txMessages = append(txMessages, msgParsed)
		}
	}

	if len(txMessages) == 0 {
		c.Logger.Debug().
			Int64("height", txResult.Height).
			Str("hash", txHash).
			Msg("All messages in transaction were filtered out, skipping.")
		return nil
	}

	return &types.Tx{
		Hash:          c.Chain.GetTransactionLink(txHash),
		Height:        c.Chain.GetBlockLink(txResult.Height),
		Memo:          txProto.GetBody().GetMemo(),
		Messages:      txMessages,
		MessagesCount: len(txProto.GetBody().GetMessages()),
		Code:          txResult.Result.Code,
		Log:           txResult.Result.Log,
	}
}

func (c *Converter) ParseMessage(
	message *codecTypes.Any,
	height int64,
) types.Message {
	parser, ok := c.Parsers[message.TypeUrl]
	if !ok {
		return &messages.MsgUnsupportedMessage{MsgType: message.TypeUrl}
	}

	msgParsed, err := parser(message.Value, c.Chain, height)
	if err != nil {
		return &messages.MsgUnparsedMessage{
			MsgType: message.TypeUrl,
			Error:   fmt.Errorf("Error parsing message: %s", err),
		}
	} else if msgParsed == nil {
		c.Logger.Error().Str("type", message.TypeUrl).Msg("Got empty message after parsing")
		return nil
	}

	// Processing internal messages (such as ones in MsgExec
	for _, internalMessage := range msgParsed.GetRawMessages() {
		if internalMessageParsed := c.ParseMessage(internalMessage, height); internalMessageParsed != nil {
			msgParsed.AddParsedMessage(internalMessageParsed)
		}
	}

	c.Logger.Debug().Str("type", message.TypeUrl).Msg("Got message")
	return msgParsed
}
