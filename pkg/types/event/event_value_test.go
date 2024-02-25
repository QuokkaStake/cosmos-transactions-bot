package event_test

import (
	eventPkg "main/pkg/types/event"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventValueFrom(t *testing.T) {
	t.Parallel()

	event := eventPkg.From("event", "key", "value")

	require.Equal(t, "event.key", event.Key)
	require.Equal(t, "value", event.Value)
}

func TestEventValueMap(t *testing.T) {
	t.Parallel()

	event1 := eventPkg.From("event", "key", "value1")
	event2 := eventPkg.From("event", "key", "value2")
	event3 := eventPkg.From("event", "otherkey", "value")

	events := eventPkg.EventValues{event1, event2, event3}
	eventsMap := events.ToMap()

	require.Len(t, eventsMap, 2)

	require.Equal(t, []string{"value1", "value2"}, eventsMap["event.key"])
	require.Equal(t, []string{"value"}, eventsMap["event.otherkey"])
}
