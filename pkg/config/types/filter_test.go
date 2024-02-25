package types_test

import (
	configTypes "main/pkg/config/types"
	"main/pkg/types/event"
	"testing"

	queryPkg "github.com/cometbft/cometbft/libs/pubsub/query"
	"github.com/stretchr/testify/require"
)

func TestFiltersString(t *testing.T) {
	t.Parallel()

	query1 := queryPkg.MustParse("event1.key1 = 'value1'")
	query2 := queryPkg.MustParse("event2.key2 = 'value2'")

	filters := configTypes.Filters{*query1, *query2}

	require.Equal(t, "event1.key1 = 'value1', event2.key2 = 'value2'", filters.String())
}

func TestFiltersMatchesEmpty(t *testing.T) {
	t.Parallel()

	filters := configTypes.Filters{}
	eventValues := event.EventValues{
		{},
	}

	matches, err := filters.Matches(eventValues)
	require.True(t, matches)
	require.NoError(t, err)
}

func TestFiltersMatches(t *testing.T) {
	t.Parallel()

	query := queryPkg.MustParse("event.key = 'value'")

	filters := configTypes.Filters{*query}
	eventValues := event.EventValues{
		{Key: "event.key", Value: "value"},
	}

	matches, err := filters.Matches(eventValues)
	require.True(t, matches)
	require.NoError(t, err)
}

func TestFiltersNotMatches(t *testing.T) {
	t.Parallel()

	query := queryPkg.MustParse("event.key = 'value'")

	filters := configTypes.Filters{*query}
	eventValues := event.EventValues{
		{Key: "event.key", Value: "value2"},
	}

	matches, err := filters.Matches(eventValues)
	require.False(t, matches)
	require.NoError(t, err)
}

func TestFiltersError(t *testing.T) {
	t.Parallel()

	query := queryPkg.MustParse("event.key > 100")

	filters := configTypes.Filters{*query}
	eventValues := event.EventValues{
		{Key: "event.key", Value: "value"},
	}

	matches, err := filters.Matches(eventValues)
	require.False(t, matches)
	require.Error(t, err)
}
