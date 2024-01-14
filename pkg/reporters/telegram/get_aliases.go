package telegram

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (reporter *TelegramReporter) HandleGetAliases(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got get aliases query")

	if !reporter.AliasManager.Enabled() {
		return reporter.BotReply(c, "Aliases manager not enabled")
	}

	aliases := reporter.AliasManager.GetAliasesLinks()
	template, err := reporter.Render("Aliases", aliases)
	if err != nil {
		return reporter.BotReply(c, fmt.Sprintf("Error displaying aliases: %s", err))
	}

	return reporter.BotReply(c, template)
}
