package telegram

import (
	"fmt"
	"main/pkg/constants"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleSetAlias(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got set alias query")

	reporter.MetricsManager.LogReporterQuery(reporter.Name(), constants.ReporterQuerySetAlias)

	if !reporter.AliasManager.Enabled() {
		return reporter.BotReply(c, "Aliases manager not enabled")
	}

	args := strings.SplitAfterN(c.Text(), " ", 4)
	if len(args) < 4 {
		return reporter.BotReply(c, fmt.Sprintf("Usage: %s <chain> <address> <alias>", args[0]))
	}

	chain, address, alias := args[1], args[2], args[3]
	chain = strings.TrimSpace(chain)
	address = strings.TrimSpace(address)
	alias = strings.TrimSpace(alias)

	chainFound := reporter.Config.Chains.FindByName(chain)
	if chainFound == nil {
		return reporter.BotReply(c, fmt.Sprintf("Chain %s is not found in config!", chain))
	}

	subscription, found := reporter.DataFetcher.FindSubscriptionByReporter(reporter.Name())
	if !found {
		return reporter.BotReply(c, "This reporter is not linked to any subscription!")
	}

	if err := reporter.AliasManager.Set(subscription.Name, chain, address, alias); err != nil {
		return reporter.BotReply(c, fmt.Sprintf("Error saving alias: %s", err))
	}

	return reporter.BotReply(c, fmt.Sprintf(
		"Saved alias: %s -> <code>%s</code>",
		reporter.SerializeLink(chainFound.GetWalletLink(address)),
		alias,
	))
}
