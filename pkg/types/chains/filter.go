package chains

import "strings"

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
		Value:    split[2],
	}
}

func (f Filter) Matches(values map[string]string) bool {
	value, found := values[f.Key]
	if !found {
		return true
	}

	switch f.Operator {
	case "=":
		return value == f.Value
	case "!=":
		return value != f.Value
	}

	// should not reach here
	panic("Received unexpected operator: " + f.Operator)
}

type Filters []Filter

func (f Filters) Matches(values map[string]string) bool {
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
