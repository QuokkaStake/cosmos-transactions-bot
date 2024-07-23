package messages

import (
	"github.com/stretchr/testify/require"
	"main/pkg/types"
	"testing"
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
