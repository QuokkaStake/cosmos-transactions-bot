package types

type Subscriptions []*Subscription

type Subscription struct {
	Name     string
	Chain    string
	Reporter string
	Filters  Filters

	LogUnknownMessages     bool
	LogUnparsedMessages    bool
	LogFailedTransactions  bool
	LogNodeErrors          bool
	FilterInternalMessages bool
}
