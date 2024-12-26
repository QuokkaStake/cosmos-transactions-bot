package telegram

import (
	"fmt"
	"main/pkg/constants"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) GetGetAliasesCommand() Command {
	return Command{
		Name:    "help",
		Query:   constants.ReporterQueryGetAliases,
		Execute: reporter.HandleGetAliases,
	}
}

func (reporter *Reporter) HandleGetAliases(c tele.Context) (string, error) {
	if !reporter.AliasManager.Enabled() {
		return "Aliases manager is not enabled!", fmt.Errorf("aliases manager not enabled")
	}

	subscription, found := reporter.DataFetcher.FindSubscriptionByReporter(reporter.Name())
	if !found {
		return "This reporter is not linked to any subscription!", fmt.Errorf("no subscriptions")
	}

	aliases := reporter.AliasManager.GetAliasesLinks(subscription.Name)
	return reporter.Render("Aliases", aliases)
}
