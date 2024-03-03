package nodes_manager_test

import (
	configTypes "main/pkg/config/types"
	nodesManagerPkg "main/pkg/nodes_manager"
	"main/pkg/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueueAdd(t *testing.T) {
	t.Parallel()

	queue := nodesManagerPkg.NewReportQueue(1)
	report1 := types.Report{Reportable: &types.Tx{Hash: configTypes.Link{Value: "123"}}}
	report2 := types.Report{Reportable: &types.Tx{Hash: configTypes.Link{Value: "456"}}}

	queue.Add(report1)
	require.True(t, queue.Has(report1))
	require.False(t, queue.Has(report2))

	queue.Add(report2)
	require.True(t, queue.Has(report2))
	require.False(t, queue.Has(report1))
}
