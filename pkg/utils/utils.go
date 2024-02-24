package utils

import (
	"strings"
)

func Map[T any, V any](source []T, mapper func(T) V) []V {
	destination := make([]V, len(source))

	for index, elt := range source {
		destination[index] = mapper(elt)
	}

	return destination
}

func Contains[T comparable](slice []T, value T) bool {
	for _, elt := range slice {
		if elt == value {
			return true
		}
	}

	return false
}

func RemoveFirstSlash(str string) string {
	if len(str) == 0 {
		return str
	}

	if str[0] == '/' {
		return str[1:]
	}

	return str
}

func SplitStringIntoChunks(msg string, maxLineLength int) []string {
	msgsByNewline := strings.Split(msg, "\n")
	outMessages := []string{}

	var sb strings.Builder

	for _, line := range msgsByNewline {
		if sb.Len()+len(line) > maxLineLength {
			outMessages = append(outMessages, sb.String())
			sb.Reset()
		}

		sb.WriteString(line + "\n")
	}

	outMessages = append(outMessages, sb.String())
	return outMessages
}

func StripTrailingDigits(s string, digits int) string {
	if i := strings.Index(s, "."); i >= 0 {
		if digits <= 0 {
			return s[:i]
		}
		i++
		if i+digits >= len(s) {
			return s
		}
		return s[:i+digits]
	}
	return s
}

func BoolToFloat64(value bool) float64 {
	if value {
		return 1
	}

	return 0
}
