package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxError(t *testing.T) {
	t.Parallel()

	event := TxError{}
	assert.NotEmpty(t, event.Type())
	assert.NotEmpty(t, event.GetHash())
	assert.Empty(t, event.GetMessages())

	event.GetAdditionalData(nil, "test")
}
