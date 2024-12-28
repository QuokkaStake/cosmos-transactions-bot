package telegram

import (
	"fmt"
	"main/pkg/constants"
	"main/pkg/types"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) GetListNodesCommand() Command {
	return Command{
		Name:    "help",
		Query:   constants.ReporterQueryNodesStatus,
		Execute: reporter.HandleListNodesStatus,
	}
}

func (reporter *Reporter) HandleListNodesStatus(c tele.Context) (string, error) {
	chains := reporter.DataFetcher.FindChainsByReporter(reporter.Name())
	if len(chains) == 0 {
		return "This reporter is not linked to any chains!", fmt.Errorf("no chains linked")
	}

	statuses := map[string]map[string]types.TendermintRPCStatus{}

	for chain, chainNodes := range reporter.NodesManager.Nodes {
		if !chains.HasChain(chain) {
			continue
		}

		statuses[chain] = map[string]types.TendermintRPCStatus{}
		for _, node := range chainNodes {
			statuses[chain][node.URL] = node.Status()
		}
	}

	return reporter.TemplatesManager.Render("Status", statuses)
}
