package telegram

import (
	"main/pkg/config/types"
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

type SetAliasRender struct {
	Alias   string
	Address string
	Chain   *types.Chain
}
