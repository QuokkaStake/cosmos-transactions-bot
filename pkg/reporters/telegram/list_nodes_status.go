package telegram

import (
	"bytes"
	tele "gopkg.in/telebot.v3"
	"main/pkg/types"
)

func (reporter *TelegramReporter) HandleListNodesStatus(c tele.Context) error {
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

	template, err := reporter.GetTemplate("Status")
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", "status").Msg("Error loading template")
		return err
	}

	var buffer bytes.Buffer
	err = template.Execute(&buffer, statuses)
	if err != nil {
		reporter.Logger.Error().Err(err).Str("type", "status").Msg("Error rendering template")
		return err
	}

	return reporter.BotReply(c, buffer.String())
}
