package nodes_manager

import (
	metricsPkg "main/pkg/metrics"
	"main/pkg/types"
	"sync"

	"main/pkg/config"
	"main/pkg/tendermint/ws"

	"github.com/rs/zerolog"
)

type NodesManager struct {
	Logger         zerolog.Logger
	Nodes          map[string][]*ws.TendermintWebsocketClient
	MetricsManager *metricsPkg.Manager

	Channel chan types.Report
	Queue   ReportQueue
	Mutex   sync.Mutex
}

func NewNodesManager(
	logger *zerolog.Logger,
	config *config.AppConfig,
	metricsManager *metricsPkg.Manager,
) *NodesManager {
	nodes := make(map[string][]*ws.TendermintWebsocketClient, len(config.Chains))

	for _, chain := range config.Chains {
		nodes[chain.Name] = make([]*ws.TendermintWebsocketClient, len(chain.TendermintNodes))

		for index, node := range chain.TendermintNodes {
			nodes[chain.Name][index] = ws.NewTendermintClient(
				logger,
				node,
				chain,
				metricsManager,
			)
		}
	}

	return &NodesManager{
		Logger:         logger.With().Str("component", "nodes_manager").Logger(),
		MetricsManager: metricsManager,
		Nodes:          nodes,
		Channel:        make(chan types.Report),
		Queue:          NewReportQueue(100),
	}
}

func (m *NodesManager) Listen() {
	for _, chain := range m.Nodes {
		for _, node := range chain {
			go node.Listen()
		}
	}

	for _, chain := range m.Nodes {
		for _, node := range chain {
			go func(c chan types.Report) {
				for msg := range c {
					m.Mutex.Lock()

					if m.Queue.Has(msg) {
						m.Logger.Trace().
							Str("hash", msg.Reportable.GetHash()).
							Msg("Message already received, not sending again.")
						m.Mutex.Unlock()
						continue
					}

					m.Channel <- msg
					m.Queue.Add(msg)

					m.Mutex.Unlock()
				}
			}(node.Channel)
		}
	}
}

func (m *NodesManager) Stop() {
	for _, chain := range m.Nodes {
		for _, node := range chain {
			node.Stop()
		}
	}
}
