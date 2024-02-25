package types

import (
	"testing"

	"github.com/rs/zerolog"
)

func TestDisplayWarningLog(t *testing.T) {
	t.Parallel()

	warning := DisplayWarning{
		Keys: map[string]string{"chain": "chain"},
	}

	logger := zerolog.Logger{}
	warning.Log(&logger)
}
