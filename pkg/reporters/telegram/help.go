package telegram

import (
	"main/pkg/constants"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) HandleHelp(c tele.Context) error {
	reporter.Logger.Info().
		Str("sender", c.Sender().Username).
		Str("text", c.Text()).
		Msg("Got help query")

	reporter.MetricsManager.LogReporterQuery(reporter.Name(), constants.ReporterQueryHelp)

	template, err := reporter.Render("Help", reporter.Version)
	if err != nil {
		return err
	}

	return reporter.BotReply(c, template)
}
