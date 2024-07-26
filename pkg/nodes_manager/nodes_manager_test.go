package nodes_manager

import (
	configPkg "main/pkg/config"
	"main/pkg/config/types"
	loggerPkg "main/pkg/logger"
	"main/pkg/metrics"
	types2 "main/pkg/types"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // disabled due to httpmock usage
func TestNodesManagerReceive(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config := &configPkg.AppConfig{
		Chains: types.Chains{
			{Name: "chain", TendermintNodes: []string{"example"}},
		},
		Metrics: configPkg.MetricsConfig{Enabled: false},
	}

	logger := loggerPkg.GetNopLogger()
	metricsManager := metrics.NewManager(logger, config.Metrics)
	nodesManager := NewNodesManager(logger, config, metricsManager)

	go nodesManager.Listen()
	defer nodesManager.Stop()

	reportable := &types2.Tx{Hash: types.Link{Value: "123"}}

	nodesManager.Nodes["chain"][0].Channel <- types2.Report{
		Chain:      config.Chains[0],
		Reportable: reportable,
	}

	received := <-nodesManager.Channel

	require.Equal(t, "123", received.Reportable.GetHash())

	nodesManager.Nodes["chain"][0].Channel <- types2.Report{
		Chain:      config.Chains[0],
		Reportable: reportable,
	}
}
