package messages

import (
	"main/pkg/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMsgUnsupportedMessageBase(t *testing.T) {
	t.Parallel()

	msg := MsgUnsupportedMessage{MsgType: "type"}

	require.Equal(t, "MsgUnsupportedMessage", msg.Type())

	msg.AddParsedMessage(nil)
	msg.SetParsedMessages([]types.Message{})
	msg.GetAdditionalData(nil, "subscription")

	require.Empty(t, msg.GetValues())
	require.Empty(t, msg.GetParsedMessages())
	require.Empty(t, msg.GetRawMessages())
}
