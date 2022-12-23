package utils

func Map[T any, V any](source []T, mapper func(T) V) []V {
	destination := make([]V, len(source))

	for index, elt := range source {
		destination[index] = mapper(elt)
	}

	return destination
}
