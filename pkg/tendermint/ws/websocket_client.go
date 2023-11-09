package ws

import (
	"context"
	metricsPkg "main/pkg/metrics"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/cometbft/cometbft/libs/pubsub/query"

	configTypes "main/pkg/config/types"
	"main/pkg/converter"
	"main/pkg/types"

	tmClient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	jsonRpcTypes "github.com/cometbft/cometbft/rpc/jsonrpc/types"
	"github.com/rs/zerolog"
)

type TendermintWebsocketClient struct {
	Logger         zerolog.Logger
	Chain          configTypes.Chain
	MetricsManager *metricsPkg.Manager
	URL            string
	Queries        []query.Query
	Client         *tmClient.WSClient
	Converter      *converter.Converter
	Active         bool
	Error          error

	Parsers map[string]types.MessageParser
	Channel chan types.Report
}

func NewTendermintClient(
	logger *zerolog.Logger,
	url string,
	chain *configTypes.Chain,
	metricsManager *metricsPkg.Manager,
) *TendermintWebsocketClient {
	return &TendermintWebsocketClient{
		Logger: logger.With().
			Str("component", "tendermint_ws_client").
			Str("url", url).
			Str("chain", chain.Name).
			Logger(),
		MetricsManager: metricsManager,
		URL:            url,
		Chain:          *chain,
		Queries:        chain.Queries,
		Active:         false,
		Channel:        make(chan types.Report),
		Converter:      converter.NewConverter(logger, chain),
	}
}

func (t *TendermintWebsocketClient) Status() types.TendermintRPCStatus {
	if t.Client == nil {
		return types.TendermintRPCStatus{
			Success: false,
			Error:   t.Error,
		}
	}

	return types.TendermintRPCStatus{
		Success: t.Active,
		Error:   t.Error,
	}
}

func SetUnexportedField(field reflect.Value, value interface{}) {
	reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}

func (t *TendermintWebsocketClient) Listen() {
	client, err := tmClient.NewWS(
		t.URL,
		"/websocket",
		tmClient.OnReconnect(func() {
			t.Logger.Info().Msg("Reconnecting...")
			t.SubscribeToUpdates()
		}),
		tmClient.PingPeriod(1*time.Second),
	)
	if err != nil {
		t.Logger.Error().Err(err).Msg("Failed to create a client")
		t.Error = err
		t.Channel <- t.MakeReport(&types.NodeConnectError{Error: err, URL: t.URL, Chain: t.Chain.GetName()})
		t.MetricsManager.LogNodeConnection(t.Chain.Name, t.URL, false)
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
		t.Channel <- t.MakeReport(&types.NodeConnectError{Error: err, URL: t.URL, Chain: t.Chain.GetName()})
		t.Logger.Warn().Err(err).Msg("Error connecting to node")
		t.MetricsManager.LogNodeConnection(t.Chain.Name, t.URL, false)
	} else {
		t.Logger.Debug().Msg("Connected to a node")
		t.Active = true
		t.MetricsManager.LogNodeConnection(t.Chain.Name, t.URL, true)
	}

	t.SubscribeToUpdates()

	for {
		select {
		case result := <-client.ResponsesCh:
			t.ProcessEvent(result)
		}
	}
}

func (t *TendermintWebsocketClient) Stop() {
	t.Logger.Info().Msg("Stopping the node...")

	if t.Client != nil {
		if err := t.Client.Stop(); err != nil {
			t.Logger.Warn().Err(err).Msg("Error stopping the node")
		}
	}
}

func (t *TendermintWebsocketClient) SubscribeToUpdates() {
	t.Logger.Trace().Msg("Subscribing to updates...")

	for _, query := range t.Queries {
		if err := t.Client.Subscribe(context.Background(), query.String()); err != nil {
			t.Logger.Error().Err(err).Str("query", query.String()).Msg("Failed to subscribe to query")
		} else {
			t.Logger.Info().Str("query", query.String()).Msg("Listening for incoming transactions")
		}
	}
}

func (t *TendermintWebsocketClient) ProcessEvent(event jsonRpcTypes.RPCResponse) {
	reportable := t.Converter.ParseEvent(event)
	if reportable != nil {
		t.MetricsManager.LogWSEvent(t.Chain.Name, t.URL)
		t.Channel <- t.MakeReport(reportable)
	}
}

func (t *TendermintWebsocketClient) MakeReport(reportable types.Reportable) types.Report {
	return types.Report{
		Chain:      t.Chain,
		Node:       t.URL,
		Reportable: reportable,
	}
}
