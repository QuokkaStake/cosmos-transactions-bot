package telegram

import (
	"fmt"
	"main/pkg/constants"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleGetAliases(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got get aliases query")

	reporter.MetricsManager.LogReporterQuery(reporter.Name(), constants.ReporterQueryGetAliases)

	if !reporter.AliasManager.Enabled() {
		return reporter.BotReply(c, "Aliases manager not enabled!")
	}

	subscription, found := reporter.DataFetcher.FindSubscriptionByReporter(reporter.Name())
	if !found {
		return reporter.BotReply(c, "This reporter is not linked to any subscription!")
	}

	aliases := reporter.AliasManager.GetAliasesLinks(subscription.Name)
	template, err := reporter.Render("Aliases", aliases)
	if err != nil {
		return reporter.BotReply(c, fmt.Sprintf("Error displaying aliases: %s", err))
	}

	return reporter.BotReply(c, template)
}
