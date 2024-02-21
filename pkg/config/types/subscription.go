package types

type Subscriptions []*Subscription

type Subscription struct {
	Name               string
	Reporter           string
	ChainSubscriptions ChainSubscriptions
}

type ChainSubscriptions []*ChainSubscription

type ChainSubscription struct {
	Chain   string
	Filters Filters

	LogUnknownMessages     bool
	LogUnparsedMessages    bool
	LogFailedTransactions  bool
	LogNodeErrors          bool
	FilterInternalMessages bool
}
