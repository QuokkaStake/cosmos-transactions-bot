package telegram

import (
	"fmt"
	"main/pkg/constants"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) GetSetAliasCommand() Command {
	return Command{
		Name:    "alias",
		Query:   constants.ReporterQuerySetAlias,
		Execute: reporter.HandleSetAlias,
		MinArgs: 3,
		Usage:   "Usage: %s <chain> <address> <alias>",
	}
}

func (reporter *Reporter) HandleSetAlias(c tele.Context) (string, error) {
	if !reporter.AliasManager.Enabled() {
		return "Aliases manager is not enabled!", fmt.Errorf("aliases manager not enabled")
	}

	args := strings.SplitAfterN(c.Text(), " ", 4)

	chain, address, alias := args[1], args[2], args[3]
	chain = strings.TrimSpace(chain)
	address = strings.TrimSpace(address)
	alias = strings.TrimSpace(alias)

	chainFound := reporter.Config.Chains.FindByName(chain)
	if chainFound == nil {
		return fmt.Sprintf("Chain %s is not found in config!", chain), fmt.Errorf("chain not found")
	}

	subscription, found := reporter.DataFetcher.FindSubscriptionByReporter(reporter.Name())
	if !found {
		return "This reporter is not linked to any subscription!", fmt.Errorf("no subscriptions")
	}

	if err := reporter.AliasManager.Set(subscription.Name, chain, address, alias); err != nil {
		return fmt.Sprintf("Error saving alias: %s", err), err
	}

	return reporter.TemplatesManager.Render("SetAlias", SetAliasRender{
		Chain:   chainFound,
		Alias:   alias,
		Address: address,
	})
}
