package types

import (
	"main/pkg/data_fetcher"
)

type Message interface {
	Type() string
	GetAdditionalData(data_fetcher.DataFetcher)
}
