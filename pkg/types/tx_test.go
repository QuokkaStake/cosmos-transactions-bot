package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTxGetMessagesLabel(t *testing.T) {
	t.Parallel()

	event := Tx{
		Messages: []Message{
			nil,
			nil,
		},
		MessagesCount: 2,
	}

	require.Equal(t, "2", event.GetMessagesLabel())

	event2 := Tx{
		Messages: []Message{
			nil,
		},
		MessagesCount: 2,
	}

	require.Equal(t, "2, 1 skipped", event2.GetMessagesLabel())
}
