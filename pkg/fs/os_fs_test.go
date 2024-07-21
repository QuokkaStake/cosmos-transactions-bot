package fs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOsFsCreate(t *testing.T) {
	t.Parallel()

	fs := &OsFS{}
	_, err := fs.Create("/tmp/file.txt")
	require.NoError(t, err)
}
