package types

import "github.com/rs/zerolog"

type DisplayWarning struct {
	Keys map[string]string
	Text string
}

func (d DisplayWarning) Log(logger *zerolog.Logger) {
	event := logger.Warn()

	for key, value := range d.Keys {
		event = event.Str(key, value)
	}

	event.Msg(d.Text)
}
