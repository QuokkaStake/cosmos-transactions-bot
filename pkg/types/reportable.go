package types

type Reportable interface {
	Type() string
	GetHash() string
	GetMessages() []Message
	GetAdditionalData(dataFetcher DataFetcher, subscriptionName string)
}
