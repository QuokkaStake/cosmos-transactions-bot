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
		if sb.Len()+len(line) >= maxLineLength {
			outMessages = append(outMessages, sb.String())
			sb.Reset()
		}

		sb.WriteString(line + "\n")
	}

	outMessages = append(outMessages, sb.String())
	return outMessages
}
