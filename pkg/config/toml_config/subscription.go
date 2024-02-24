package toml_config

import (
	"fmt"
	"main/pkg/config/types"

	"github.com/cometbft/cometbft/libs/pubsub/query"
	"gopkg.in/guregu/null.v4"
)

type Subscriptions []*Subscription

type Subscription struct {
	Name              string             `toml:"name"`
	Reporter          string             `toml:"reporter"`
	ChainSubscription ChainSubscriptions `toml:"chains"`
}

type ChainSubscriptions []*ChainSubscription

type ChainSubscription struct {
	Chain                  string    `toml:"name"`
	Filters                []string  `toml:"filters"`
	LogUnknownMessages     null.Bool `default:"false" toml:"log-unknown-messages"`
	LogUnparsedMessages    null.Bool `default:"true"  toml:"log-unparsed-messages"`
	LogFailedTransactions  null.Bool `default:"true"  toml:"log-failed-transactions"`
	LogNodeErrors          null.Bool `default:"true"  toml:"log-node-errors"`
	FilterInternalMessages null.Bool `default:"true"  toml:"filter-internal-messages"`
}

func (subscriptions Subscriptions) Validate() error {
	for index, subscription := range subscriptions {
		if err := subscription.Validate(); err != nil {
			return fmt.Errorf("error in subscription %d: %s", index, err)
		}
	}

	// checking names uniqueness
	names := map[string]bool{}

	for _, subscription := range subscriptions {
		if _, ok := names[subscription.Name]; ok {
			return fmt.Errorf("duplicate subscription name: %s", subscription.Name)
		}

		names[subscription.Name] = true
	}

	return nil
}

func (s *Subscription) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("empty subscription name")
	}

	if s.Reporter == "" {
		return fmt.Errorf("empty reporter name")
	}

	for index, subscription := range s.ChainSubscription {
		if err := subscription.Validate(); err != nil {
			return fmt.Errorf("error in subscription %d: %s", index, err)
		}
	}

	return nil
}

func (s *ChainSubscription) Validate() error {
	if s.Chain == "" {
		return fmt.Errorf("empty chain name")
	}

	for index, filter := range s.Filters {
		if _, err := query.New(filter); err != nil {
			return fmt.Errorf("error in filter %d: %s", index, err)
		}
	}

	return nil
}

func (s *ChainSubscription) ToAppConfigChainSubscription() *types.ChainSubscription {
	filters := make([]query.Query, len(s.Filters))
	for index, filter := range s.Filters {
		filters[index] = *query.MustParse(filter)
	}

	return &types.ChainSubscription{
		Chain:                  s.Chain,
		Filters:                filters,
		LogUnknownMessages:     s.LogUnknownMessages.Bool,
		LogUnparsedMessages:    s.LogUnparsedMessages.Bool,
		LogFailedTransactions:  s.LogFailedTransactions.Bool,
		LogNodeErrors:          s.LogNodeErrors.Bool,
		FilterInternalMessages: s.FilterInternalMessages.Bool,
	}
}

func (s *Subscription) ToAppConfigSubscription() *types.Subscription {
	chainSubscriptions := make(types.ChainSubscriptions, len(s.ChainSubscription))
	for index, chainSubscription := range s.ChainSubscription {
		chainSubscriptions[index] = chainSubscription.ToAppConfigChainSubscription()
	}

	return &types.Subscription{
		Name:               s.Name,
		Reporter:           s.Reporter,
		ChainSubscriptions: chainSubscriptions,
	}
}

func FromAppConfigChainSubscription(s *types.ChainSubscription) *ChainSubscription {
	subscription := &ChainSubscription{
		Chain:                  s.Chain,
		LogUnknownMessages:     null.BoolFrom(s.LogUnknownMessages),
		LogUnparsedMessages:    null.BoolFrom(s.LogUnparsedMessages),
		LogFailedTransactions:  null.BoolFrom(s.LogFailedTransactions),
		LogNodeErrors:          null.BoolFrom(s.LogNodeErrors),
		FilterInternalMessages: null.BoolFrom(s.FilterInternalMessages),
	}

	subscription.Filters = make([]string, len(s.Filters))
	for index, filter := range s.Filters {
		subscription.Filters[index] = filter.String()
	}

	return subscription
}

func FromAppConfigSubscription(s *types.Subscription) *Subscription {
	subscription := &Subscription{
		Name:              s.Name,
		Reporter:          s.Reporter,
		ChainSubscription: make(ChainSubscriptions, len(s.ChainSubscriptions)),
	}

	for index, chainSubscription := range s.ChainSubscriptions {
		subscription.ChainSubscription[index] = FromAppConfigChainSubscription(chainSubscription)
	}

	return subscription
}
