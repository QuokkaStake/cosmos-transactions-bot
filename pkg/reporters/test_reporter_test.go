package reporters

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTestReporter(t *testing.T) {
	t.Parallel()

	reporter := &TestReporter{ReporterName: "test"}
	reporter.Init()
	assert.Equal(t, "test", reporter.Name())
	assert.Equal(t, "telegram", reporter.Type())
}
