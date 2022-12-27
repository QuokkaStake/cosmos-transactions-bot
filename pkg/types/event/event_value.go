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
