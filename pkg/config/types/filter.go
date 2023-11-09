package types

import (
	"main/pkg/types/event"
	"strings"

	"github.com/cometbft/cometbft/libs/pubsub/query"
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
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

func (f Filters) MatchesType(msgType string) (bool, error) {
	values := []event.EventValue{
		event.From(cosmosTypes.EventTypeMessage, cosmosTypes.AttributeKeyAction, msgType),
	}

	return f.Matches(values)
}
