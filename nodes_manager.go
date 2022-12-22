package main

import (
	"sync"

	"github.com/rs/zerolog"
)

type ReportQueue struct {
	Data  []Report
	Size  int
	Mutes sync.Mutex
}

func NewReportQueue(size int) ReportQueue {
	return ReportQueue{Data: make([]Report, 0), Size: size}
}

func (q *ReportQueue) Add(report Report) {
	q.Mutes.Lock()

	if len(q.Data) >= q.Size {
		_, q.Data = q.Data[0], q.Data[1:]
	}

	q.Data = append(q.Data, report)
	q.Mutes.Unlock()
}

func (q *ReportQueue) Has(msg Report) bool {
	for _, elem := range q.Data {
		if elem.Reportable.GetHash() == msg.Reportable.GetHash() {
			return true
		}
	}

	return false
}

type NodesManager struct {
	Logger  zerolog.Logger
	Nodes   map[string][]*TendermintClient
	Channel chan Report
	Queue   ReportQueue
	Mutex   sync.Mutex
}

func NewNodesManager(logger *zerolog.Logger, config *Config) *NodesManager {
	nodes := make(map[string][]*TendermintClient, len(config.Chains))

	for _, chain := range config.Chains {
		nodes[chain.Name] = make([]*TendermintClient, len(chain.TendermintNodes))

		for index, node := range chain.TendermintNodes {
			nodes[chain.Name][index] = NewTendermintClient(
				logger,
				node,
				chain,
			)
		}
	}

	return &NodesManager{
		Logger:  logger.With().Str("component", "nodes_manager").Logger(),
		Nodes:   nodes,
		Channel: make(chan Report),
		Queue:   NewReportQueue(100),
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
			go func(c chan Report) {
				for msg := range c {
					m.Mutex.Lock()

					if m.Queue.Has(msg) {
						m.Logger.Trace().
							Str("hash", msg.Reportable.GetHash()).
							Msg("Message already received, not sending again.")
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
