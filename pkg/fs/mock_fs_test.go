package fs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockFsWriteError(t *testing.T) {
	t.Parallel()

	fs := &MockFs{}
	_, err := fs.Write([]byte{})
	require.Error(t, err)
}
