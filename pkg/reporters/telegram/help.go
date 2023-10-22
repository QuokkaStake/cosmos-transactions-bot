package telegram

import (
	tele "gopkg.in/telebot.v3"
)

func (reporter *TelegramReporter) HandleHelp(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got help query")

	template, err := reporter.Render("Help", reporter.Version)
	if err != nil {
		return err
	}

	return reporter.BotReply(c, template)
}
