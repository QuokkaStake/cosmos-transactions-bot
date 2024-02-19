package telegram

import (
	"main/pkg/types"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleListNodesStatus(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got status query")

	statuses := map[string]map[string]types.TendermintRPCStatus{}

	for chain, chainNodes := range reporter.NodesManager.Nodes {
		statuses[chain] = map[string]types.TendermintRPCStatus{}
		for _, node := range chainNodes {
			statuses[chain][node.URL] = node.Status()
		}
	}

	template, err := reporter.Render("Status", statuses)
	if err != nil {
		return err
	}

	return reporter.BotReply(c, template)
}
