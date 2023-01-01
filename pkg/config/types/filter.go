package types

import (
	"fmt"
	"strings"

	"main/pkg/types/event"
	"main/pkg/utils"
)

type Filter struct {
	Key      string
	Operator string
	Value    string
}

func NewFilter(filter string) Filter {
	split := strings.Split(filter, " ")

	return Filter{
		Key:      split[0],
		Operator: split[1],
		Value:    utils.Dequotify(split[2]),
	}
}

func (f Filter) ToString() string {
	return fmt.Sprintf("%s %s '%s'", f.Key, f.Operator, f.Value)
}

func (f Filter) Matches(values event.EventValues) bool {
	for _, value := range values {
		if value.Key != f.Key {
			continue
		}

		if f.Operator == "=" && f.Value != value.Value {
			return false
		}
		if f.Operator == "!=" && f.Value == value.Value {
			return false
		}
	}

	return true
}

type Filters []Filter

func (f Filters) Matches(values event.EventValues) bool {
	if len(f) == 0 {
		return true
	}

	for _, filter := range f {
		if filter.Matches(values) {
			return true
		}
	}

	return false
}
