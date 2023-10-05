package types

import (
	"main/pkg/data_fetcher"
)

type Reportable interface {
	Type() string
	GetHash() string
	GetMessages() []Message
	GetAdditionalData(data_fetcher.DataFetcher)
}
