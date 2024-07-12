package responses_test

import (
	"encoding/json"
	"main/pkg/types/responses"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalDurationInvalidInput(t *testing.T) {
	t.Parallel()

	duration := responses.Duration{}
	err := duration.UnmarshalJSON([]byte{})
	require.Error(t, err)
}

func TestUnmarshalDurationInvalidString(t *testing.T) {
	t.Parallel()

	var duration responses.Duration
	err := json.Unmarshal([]byte("\"invalid\""), &duration)

	require.Error(t, err)
	require.Zero(t, duration.Duration)
}

func TestUnmarshalDurationNotString(t *testing.T) {
	t.Parallel()

	var duration responses.Duration
	err := json.Unmarshal([]byte("3"), &duration)

	require.Error(t, err)
	require.Zero(t, duration.Duration)
}

func TestUnmarshalDurationValid(t *testing.T) {
	t.Parallel()

	var duration responses.Duration
	err := json.Unmarshal([]byte("\"20s\""), &duration)

	require.NoError(t, err)
	require.Equal(t, 20*time.Second, duration.Duration)
}

func TestUnMarshalDurationValid(t *testing.T) {
	t.Parallel()

	var duration responses.Duration
	err := json.Unmarshal([]byte("\"20s\""), &duration)

	require.NoError(t, err)
	require.Equal(t, 20*time.Second, duration.Duration)
}
