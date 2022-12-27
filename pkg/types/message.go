package types

import (
	dataFetcher "main/pkg/data_fetcher"
	"main/pkg/types/event"
)

type Message interface {
	Type() string
	GetAdditionalData(dataFetcher.DataFetcher)
	GetValues() event.EventValues
}
