package telegram

import (
	"main/pkg/constants"

	tele "gopkg.in/telebot.v3"
)

func (reporter *Reporter) GetHelpCommand() Command {
	return Command{
		Name:    "help",
		Query:   constants.ReporterQueryHelp,
		Execute: reporter.HandleHelp,
	}
}

func (reporter *Reporter) HandleHelp(c tele.Context) (string, error) {
	return reporter.Render("Help", reporter.Version)
}
