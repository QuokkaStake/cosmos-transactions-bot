package telegram

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (reporter *TelegramReporter) HandleGetConfig(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got get config query")

	configString, err := reporter.Config.GetConfigAsString()
	if err != nil {
		return reporter.BotReply(c, fmt.Sprintf("Error converting config to string: %s", err))
	}

	return reporter.BotReply(c, configString)
}
