package main

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/tendermint/tendermint/crypto/tmhash"

	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/gogo/protobuf/proto"
	"github.com/rs/zerolog"
	"github.com/tendermint/tendermint/libs/json"
	coreTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmClient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	jsonRpcTypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	tendermintTypes "github.com/tendermint/tendermint/types"
)

type TendermintClient struct {
	Logger  zerolog.Logger
	Chain   Chain
	URL     string
	Filters []string
	Client  *tmClient.WSClient
	Active  bool
	Error   error

	Parsers map[string]MessageParser
	Channel chan Report
}

func NewTendermintClient(
	logger *zerolog.Logger,
	url string,
	chain *Chain,
) *TendermintClient {
	return &TendermintClient{
		Logger: logger.With().
			Str("component", "tendermint_client").
			Str("url", url).
			Str("chain", chain.Name).
			Logger(),
		URL:     url,
		Chain:   *chain,
		Filters: chain.Filters,
		Active:  false,
		Channel: make(chan Report),
		Parsers: map[string]MessageParser{
			"/cosmos.bank.v1beta1.MsgSend": func(data []byte, c Chain) (Message, error) {
				return ParseMsgSend(data, chain)
			},
			"/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward": func(data []byte, c Chain) (Message, error) {
				return ParseMsgWithdrawDelegatorReward(data, chain)
			},
		},
	}
}

func (t *TendermintClient) Status() TendermintRpcStatus {
	if t.Client == nil {
		return TendermintRpcStatus{
			Success: false,
			Error:   fmt.Errorf("Tendermint RPC not initialized"),
		}
	}

	return TendermintRpcStatus{
		Success: t.Active,
		Error:   t.Error,
	}
}

func SetUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}

func (t *TendermintClient) Listen() {
	client, err := tmClient.NewWS(
		t.URL,
		"/websocket",
		tmClient.OnReconnect(func() {
			t.Logger.Info().Msg("Reconnecting...")
			t.SubscribeToUpdates()
		}),
		tmClient.PingPeriod(5*time.Second),
	)

	if err != nil {
		t.Logger.Error().Err(err).Msg("Failed to create a client")
		return
	}

	// Patching WSS connections
	if strings.HasPrefix(t.URL, "https") {
		field := reflect.ValueOf(client).Elem().FieldByName("protocol")
		SetUnexportedField(field, "wss")
	}

	t.Client = client

	t.Logger.Trace().Msg("Connecting to a node...")

	if err = client.Start(); err != nil {
		t.Error = err
		t.Logger.Warn().Err(err).Msg("Error connecting to node")
	} else {
		t.Logger.Debug().Msg("Connected to a node")
		t.Active = true
	}

	t.SubscribeToUpdates()

	for {
		select {
		case result := <-client.ResponsesCh:
			t.ProcessEvent(result)
		}
	}
}

func (t *TendermintClient) Stop() {
	t.Logger.Info().Msg("Stopping the node...")

	if t.Client != nil {
		if err := t.Client.Stop(); err != nil {
			t.Logger.Warn().Err(err).Msg("Error stopping the node")
		}
	}
}

func (t *TendermintClient) SubscribeToUpdates() {
	t.Logger.Trace().Msg("Subscribing to updates...")

	for _, filter := range t.Filters {
		if err := t.Client.Subscribe(context.Background(), filter); err != nil {
			t.Logger.Error().Err(err).Str("filter", filter).Msg("Failed to subscribe to filter")
		} else {
			t.Logger.Info().Str("filter", filter).Msg("Listening for incoming transactions")
		}
	}
}

func (t *TendermintClient) ProcessEvent(event jsonRpcTypes.RPCResponse) {
	if event.Error != nil && event.Error.Message != "" {
		t.Logger.Error().Str("msg", event.Error.Error()).Msg("Got error in RPC response")
		t.Channel <- t.MakeReport(TxError{Error: event.Error})
		return
	}

	var resultEvent coreTypes.ResultEvent
	if err := json.Unmarshal(event.Result, &resultEvent); err != nil {
		t.Logger.Error().Err(err).Msg("Failed to parse event")
		t.Channel <- t.MakeReport(TxError{Error: event.Error})
		return
	}

	if resultEvent.Data == nil {
		t.Logger.Debug().Msg("Event does not have data, skipping.")
		return
	}

	txResult := resultEvent.Data.(tendermintTypes.EventDataTx).TxResult
	txHash := fmt.Sprintf("%X", tmhash.Sum(txResult.Tx))
	var txProto tx.Tx

	if err := proto.Unmarshal(txResult.Tx, &txProto); err != nil {
		t.Logger.Error().Err(err).Msg("Could not parse tx")
		t.Channel <- t.MakeReport(TxError{Error: event.Error})
		return
	}

	t.Logger.Debug().
		Int64("height", txResult.Height).
		Str("memo", txProto.GetBody().GetMemo()).
		Str("hash", txHash).
		Int("len", len(txProto.GetBody().Messages)).
		Msg("Got transaction")

	messages := []Message{}

	for _, message := range txProto.GetBody().Messages {
		t.Logger.Debug().Str("type", message.TypeUrl).Msg("Got message")

		var msgParsed Message
		var err error

		if parser, ok := t.Parsers[message.TypeUrl]; ok {
			msgParsed, err = parser(message.Value, t.Chain)
			if err != nil {
				t.Logger.Error().Err(err).Str("type", message.TypeUrl).Msg("Error parsing message")
				msgParsed = MsgError{
					Error: fmt.Errorf("Error parsing message: %s", err),
				}
			}
		} else {
			msgParsed = MsgError{
				Error: fmt.Errorf("Got unsupported message type: %s", message.TypeUrl),
			}
		}

		if msgParsed != nil {
			messages = append(messages, msgParsed)
		}
	}

	txParsed := Tx{
		Hash:     txHash,
		Height:   txResult.Height,
		Memo:     txProto.GetBody().GetMemo(),
		Messages: messages,
	}

	t.Channel <- t.MakeReport(txParsed)
}

func (t *TendermintClient) MakeReport(reportable Reportable) Report {
	return Report{
		Chain:      t.Chain,
		Node:       t.URL,
		Reportable: reportable,
	}
}
