package reporters

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestTestReporter(t *testing.T) {
	t.Parallel()

	reporter := &TestReporter{ReporterName: "test"}
	require.NoError(t, reporter.Init())
	assert.Equal(t, "test", reporter.Name())
	assert.Equal(t, "telegram", reporter.Type())
}
