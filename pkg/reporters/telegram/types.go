package telegram

import (
	"main/pkg/constants"

	tele "gopkg.in/telebot.v3"
)

type Command struct {
	Name    string
	MinArgs int
	Usage   string
	Query   constants.ReporterQuery
	Execute func(c tele.Context) (string, error)
}
