package types_test

import (
	"main/pkg/config/types"
	"testing"

	"github.com/rs/zerolog"
)

func TestDisplayWarningLog(t *testing.T) {
	t.Parallel()

	warning := types.DisplayWarning{
		Keys: map[string]string{"chain": "chain"},
	}

	logger := zerolog.Logger{}
	warning.Log(&logger)
}
