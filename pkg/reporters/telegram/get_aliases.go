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

	aliases, err := reporter.AliasManager.GetAsToml()
	if err != nil {
		return reporter.BotReply(c, fmt.Sprintf("Error getting aliases: %s", err))
	}

	return reporter.BotReply(c, fmt.Sprintf("<pre>%s</pre>", aliases))
}
