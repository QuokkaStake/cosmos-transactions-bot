package messages

import (
	"main/pkg/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgNotExistingMessageBase(t *testing.T) {
	t.Parallel()

	msg := MsgNotExistingMessage{}

	require.Equal(t, "MsgNotExistingMessage", msg.Type())

	msg.AddParsedMessage(nil)
	msg.SetParsedMessages([]types.Message{})
	msg.GetAdditionalData(nil, "subscription")

	require.Empty(t, msg.GetValues())
	require.Empty(t, msg.GetParsedMessages())
	require.Empty(t, msg.GetRawMessages())
}
