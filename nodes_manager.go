package main

import "github.com/rs/zerolog"

type NodesManager struct {
	Logger  zerolog.Logger
	Nodes   map[string][]*TendermintClient
	Channel chan Report
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
					m.Channel <- msg
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
