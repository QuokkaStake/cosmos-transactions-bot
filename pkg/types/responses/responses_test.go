package responses_test

import (
	"encoding/json"
	"main/pkg/types/responses"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMarshalDurationInvalidString(t *testing.T) {
	t.Parallel()

	var duration responses.Duration
	err := json.Unmarshal([]byte("\"invalid\""), &duration)

	require.Error(t, err)
	require.Zero(t, duration.Duration)
}

func TestMarshalDurationNotString(t *testing.T) {
	t.Parallel()

	var duration responses.Duration
	err := json.Unmarshal([]byte("invalid"), &duration)

	require.Error(t, err)
	require.Zero(t, duration.Duration)
}

func TestMarshalDurationValid(t *testing.T) {
	t.Parallel()

	var duration responses.Duration
	err := json.Unmarshal([]byte("\"20s\""), &duration)

	require.NoError(t, err)
	require.Equal(t, 20*time.Second, duration.Duration)
}
