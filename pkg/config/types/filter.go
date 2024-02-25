package types

import (
	"main/pkg/types/event"
	"strings"

	"github.com/cometbft/cometbft/libs/pubsub/query"
)

type Filters []query.Query

func (f Filters) String() string {
	outStrings := make([]string, len(f))

	for index, filter := range f {
		outStrings[index] = filter.String()
	}

	return strings.Join(outStrings, ", ")
}

func (f Filters) Matches(values event.EventValues) (bool, error) {
	if len(f) == 0 {
		return true, nil
	}

	for _, filter := range f {
		if matches, err := filter.Matches(values.ToMap()); err != nil {
			return false, err
		} else if matches {
			return true, nil
		}
	}

	return false, nil
}
