package utils

import "strconv"

func Map[T any, V any](source []T, mapper func(T) V) []V {
	destination := make([]V, len(source))

	for index, elt := range source {
		destination[index] = mapper(elt)
	}

	return destination
}

func Contains[T comparable](array []T, element T) bool {
	for _, a := range array {
		if a == element {
			return true
		}
	}
	return false
}

func StrToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return f
}
