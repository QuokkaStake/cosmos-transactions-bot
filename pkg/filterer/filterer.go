package filterer

import (
	"fmt"
	configTypes "main/pkg/config/types"
	"main/pkg/messages"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type Filterer struct {
	Logger zerolog.Logger
	Chain  *configTypes.Chain
}

func NewFilterer(
	logger *zerolog.Logger,
	chain *configTypes.Chain,
) *Filterer {
	return &Filterer{
		Logger: logger.With().
			Str("component", "filterer").
			Str("chain", chain.Name).
			Logger(),
		Chain: chain,
	}
}

func (f *Filterer) Filter(reportable types.Reportable) types.Reportable {
	// Not filtering out TxError
	if _, ok := reportable.(*types.TxError); ok {
		return reportable
	}

	tx, ok := reportable.(*types.Tx)
	if !ok {
		f.Logger.Error().Str("type", reportable.Type()).Msg("Unsupported reportable type, ignoring.")
		return nil
	}

	if !f.Chain.LogFailedTransactions && tx.Code > 0 {
		f.Logger.Debug().
			Str("hash", tx.GetHash()).
			Msg("Transaction is failed, skipping")
		return nil
	}

	messages := make([]types.Message, 0)

	for _, message := range tx.Messages {
		filteredMessage := f.FilterMessage(message, false)
		if filteredMessage != nil {
			messages = append(messages, filteredMessage)
		}
	}

	if len(messages) == 0 {
		f.Logger.Debug().
			Str("hash", tx.GetHash()).
			Msg("All messages in transaction were filtered out, skipping.")
		return nil
	}

	tx.Messages = messages
	return tx
}

func (f *Filterer) FilterMessage(message types.Message, internal bool) types.Message {
	if unsupportedMsg, ok := message.(*messages.MsgUnsupportedMessage); ok {
		if f.Chain.LogUnknownMessages {
			f.Logger.Error().Str("type", unsupportedMsg.MsgType).Msg("Unsupported message type")
			return message
		} else {
			f.Logger.Debug().Str("type", unsupportedMsg.MsgType).Msg("Unsupported message type")
			return nil
		}
	}

	// internal -> filter only if f.Chain.FilterInternalMessages is true
	// !internal -> filter regardless
	if !internal || f.Chain.FilterInternalMessages {
		matches, err := f.Chain.Filters.Matches(message.GetValues())

		f.Logger.Trace().
			Str("type", message.Type()).
			Str("values", fmt.Sprintf("%+v", message.GetValues().ToMap())).
			Str("filters", fmt.Sprintf("%+v", f.Chain.Filters)).
			Bool("matches", matches).
			Msg("Result of matching message events against filters")

		if err != nil {
			f.Logger.Error().
				Err(err).
				Str("type", message.Type()).
				Msg("Error checking if message matches filters")
		} else if !matches {
			f.Logger.Debug().
				Str("type", message.Type()).
				Msg("Message is ignored by filters.")
			return nil
		}
	}

	if len(message.GetRawMessages()) == 0 {
		return message
	}

	parsedInternalMessages := make([]types.Message, 0)

	// Processing internal messages (such as ones in MsgExec)
	for _, internalMessage := range message.GetParsedMessages() {
		if internalMessageParsed := f.FilterMessage(internalMessage, true); internalMessageParsed != nil {
			parsedInternalMessages = append(parsedInternalMessages, internalMessageParsed)
		}
	}

	if len(parsedInternalMessages) == 0 {
		f.Logger.Debug().
			Str("type", message.Type()).
			Msg("Message with messages inside has 0 messages after filtering, skipping.")
		return nil
	}

	message.SetParsedMessages(parsedInternalMessages)
	return message
}
