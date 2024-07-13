package types

import (
	"github.com/google/uuid"
)

type UnsupportedReportable struct{}

func (e UnsupportedReportable) GetMessages() []Message {
	return []Message{}
}

func (e UnsupportedReportable) Type() string {
	return "UnsupportedReportable"
}

func (e UnsupportedReportable) GetHash() string {
	return uuid.NewString()
}

func (e *UnsupportedReportable) GetAdditionalData(fetcher DataFetcher, subscriptionName string) {
}
