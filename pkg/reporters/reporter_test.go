package reporters

import (
	configTypes "main/pkg/config/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetReporterFail(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	GetReporter(&configTypes.Reporter{}, nil, nil, nil, nil, nil, nil, "1.2.3")
}

func TestFindReporterByName(t *testing.T) {
	t.Parallel()

	reporters := Reporters{
		&TestReporter{ReporterName: "reporter"},
	}

	require.NotNil(t, reporters.FindByName("reporter"))
	require.Nil(t, reporters.FindByName("reporter2"))
}
