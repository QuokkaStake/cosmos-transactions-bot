package event

type EventValue struct {
	Key   string
	Value string
}

func From(namespace, key, value string) EventValue {
	return EventValue{
		Key:   namespace + "." + key,
		Value: value,
	}
}

type EventValues []EventValue

func (v EventValues) ToMap() map[string][]string {
	eventsMap := make(map[string][]string, len(v))

	for _, value := range v {
		if _, ok := eventsMap[value.Key]; !ok {
			eventsMap[value.Key] = []string{value.Value}
		} else {
			eventsMap[value.Key] = append(eventsMap[value.Key], value.Value)
		}
	}

	return eventsMap
}
