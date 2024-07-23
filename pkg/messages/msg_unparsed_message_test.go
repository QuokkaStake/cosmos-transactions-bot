package messages

import (
	"errors"
	"github.com/stretchr/testify/require"
	"main/pkg/types"
	"testing"
)

func TestMsgUnparsedMessageBase(t *testing.T) {
	t.Parallel()

	msg := MsgUnparsedMessage{MsgType: "type", Error: errors.New("error")}

	require.Equal(t, "MsgUnparsedMessage", msg.Type())

	msg.AddParsedMessage(nil)
	msg.SetParsedMessages([]types.Message{})
	msg.GetAdditionalData(nil, "subscription")

	require.Empty(t, msg.GetValues())
	require.Empty(t, msg.GetParsedMessages())
	require.Empty(t, msg.GetRawMessages())
}
