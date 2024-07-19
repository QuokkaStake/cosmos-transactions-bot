package metrics

import (
	"io"
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
	loggerPkg "main/pkg/logger"
	"main/pkg/messages"
	"main/pkg/types"
	"net/http"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsManagerLogLastHeight(t *testing.T) {
	t.Parallel()

	config := configPkg.MetricsConfig{Enabled: false}
	logger := loggerPkg.GetNopLogger()
	metricsManager := NewManager(logger, config)

	blockTime := time.Now()

	metricsManager.LogLastHeight("chain", 123, blockTime)

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.lastBlockHeightCollector))
	assert.InDelta(t, 123, testutil.ToFloat64(metricsManager.lastBlockHeightCollector.With(prometheus.Labels{
		"chain": "chain",
	})), 0.01)

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.lastBlockTimeCollector))
	assert.InDelta(t, blockTime.Unix(), testutil.ToFloat64(metricsManager.lastBlockTimeCollector.With(prometheus.Labels{
		"chain": "chain",
	})), 0.01)
}

func TestMetricsManagerLogNodeConnection(t *testing.T) {
	t.Parallel()

	config := configPkg.MetricsConfig{Enabled: false}
	logger := loggerPkg.GetNopLogger()
	metricsManager := NewManager(logger, config)

	metricsManager.LogNodeConnection("chain", "node1", true)
	metricsManager.LogNodeConnection("chain", "node2", false)

	assert.Equal(t, 2, testutil.CollectAndCount(metricsManager.nodeConnectedCollector))
	assert.InDelta(t, 1, testutil.ToFloat64(metricsManager.nodeConnectedCollector.With(prometheus.Labels{
		"chain": "chain",
		"node":  "node1",
	})), 0.01)
	assert.Zero(t, testutil.ToFloat64(metricsManager.nodeConnectedCollector.With(prometheus.Labels{
		"chain": "chain",
		"node":  "node2",
	})))
}

func TestMetricsManagerLogAppVersion(t *testing.T) {
	t.Parallel()

	config := configPkg.MetricsConfig{Enabled: false}
	logger := loggerPkg.GetNopLogger()
	metricsManager := NewManager(logger, config)

	metricsManager.LogAppVersion("1.2.3")

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.appVersionGauge))
	assert.InDelta(t, 1, testutil.ToFloat64(metricsManager.appVersionGauge.With(prometheus.Labels{
		"version": "1.2.3",
	})), 0.01)
}

func TestMetricsManagerLogWsEvent(t *testing.T) {
	t.Parallel()

	config := configPkg.MetricsConfig{Enabled: false}
	logger := loggerPkg.GetNopLogger()
	metricsManager := NewManager(logger, config)

	metricsManager.LogWSEvent("chain", "node")

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.eventsTotalCounter))
	assert.InDelta(t, 1, testutil.ToFloat64(metricsManager.eventsTotalCounter.With(prometheus.Labels{
		"chain": "chain",
		"node":  "node",
	})), 0.01)

	metricsManager.LogWSEvent("chain", "node")

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.eventsTotalCounter))
	assert.InDelta(t, 2, testutil.ToFloat64(metricsManager.eventsTotalCounter.With(prometheus.Labels{
		"chain": "chain",
		"node":  "node",
	})), 0.01)
}

func TestMetricsManagerLogReporterQuery(t *testing.T) {
	t.Parallel()

	config := configPkg.MetricsConfig{Enabled: false}
	logger := loggerPkg.GetNopLogger()
	metricsManager := NewManager(logger, config)

	metricsManager.LogReporterQuery("reporter", "query")

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.reporterQueriesCounter))
	assert.InDelta(t, 1, testutil.ToFloat64(metricsManager.reporterQueriesCounter.With(prometheus.Labels{
		"reporter": "reporter",
		"type":     "query",
	})), 0.01)
}

func TestMetricsManagerLogNodeReconnect(t *testing.T) {
	t.Parallel()

	config := configPkg.MetricsConfig{Enabled: false}
	logger := loggerPkg.GetNopLogger()
	metricsManager := NewManager(logger, config)

	metricsManager.LogNodeReconnect("chain", "node")

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.reconnectsCounter))
	assert.InDelta(t, 1, testutil.ToFloat64(metricsManager.reconnectsCounter.With(prometheus.Labels{
		"chain": "chain",
		"node":  "node",
	})), 0.01)

	metricsManager.LogNodeReconnect("chain", "node")

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.reconnectsCounter))
	assert.InDelta(t, 2, testutil.ToFloat64(metricsManager.reconnectsCounter.With(prometheus.Labels{
		"chain": "chain",
		"node":  "node",
	})), 0.01)
}

func TestMetricsManagerLogReporterEnabled(t *testing.T) {
	t.Parallel()

	config := configPkg.MetricsConfig{Enabled: false}
	logger := loggerPkg.GetNopLogger()
	metricsManager := NewManager(logger, config)

	metricsManager.LogReporterEnabled("name", "type")

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.reporterEnabledGauge))
	assert.InDelta(t, 1, testutil.ToFloat64(metricsManager.reporterEnabledGauge.With(prometheus.Labels{
		"name": "name",
		"type": "type",
	})), 0.01)
}

func TestMetricsManagerLogReportUnsuccessful(t *testing.T) {
	t.Parallel()

	config := configPkg.MetricsConfig{Enabled: false}
	logger := loggerPkg.GetNopLogger()
	metricsManager := NewManager(logger, config)

	report := types.Report{
		Chain:        &configTypes.Chain{Name: "chain"},
		Reportable:   &types.NodeConnectError{},
		Subscription: &configTypes.Subscription{Name: "subscription"},
	}

	metricsManager.LogReport(report, "reporter", false)

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.reporterErrorsCounter))
	assert.InDelta(t, 1, testutil.ToFloat64(metricsManager.reporterErrorsCounter.With(prometheus.Labels{
		"chain":        "chain",
		"reporter":     "reporter",
		"subscription": "subscription",
		"type":         "NodeConnectError",
	})), 0.01)
}

func TestMetricsManagerLogReportSuccessful(t *testing.T) {
	t.Parallel()

	config := configPkg.MetricsConfig{Enabled: false}
	logger := loggerPkg.GetNopLogger()
	metricsManager := NewManager(logger, config)

	report := types.Report{
		Chain: &configTypes.Chain{Name: "chain"},
		Reportable: &types.Tx{
			Messages: []types.Message{
				&messages.MsgSend{},
			},
		},
		Subscription: &configTypes.Subscription{Name: "subscription"},
	}

	metricsManager.LogReport(report, "reporter", true)

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.reporterReportsCounter))
	assert.InDelta(t, 1, testutil.ToFloat64(metricsManager.reporterReportsCounter.With(prometheus.Labels{
		"chain":        "chain",
		"reporter":     "reporter",
		"subscription": "subscription",
		"type":         "Tx",
	})), 0.01)

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.reportEntriesCounter))
	assert.InDelta(t, 1, testutil.ToFloat64(metricsManager.reportEntriesCounter.With(prometheus.Labels{
		"chain":        "chain",
		"reporter":     "reporter",
		"subscription": "subscription",
		"type":         "/cosmos.bank.v1beta1.MsgSend",
	})), 0.01)
}

func TestMetricsManagerSetDefaultMetrics(t *testing.T) {
	t.Parallel()

	config := configPkg.MetricsConfig{Enabled: false}
	logger := loggerPkg.GetNopLogger()
	metricsManager := NewManager(logger, config)

	appConfig := &configPkg.AppConfig{
		Chains: configTypes.Chains{
			{
				Name:            "chain",
				PrettyName:      "Chain Name",
				TendermintNodes: []string{"node"},
			},
		},
		Subscriptions: configTypes.Subscriptions{
			{
				Name:               "subscription",
				Reporter:           "reporter",
				ChainSubscriptions: configTypes.ChainSubscriptions{},
			},
		},
	}

	metricsManager.SetAllDefaultMetrics(appConfig)

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.startTimeGauge))

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.subscriptionsInfoCounter))
	assert.Zero(t, testutil.ToFloat64(metricsManager.subscriptionsInfoCounter.With(prometheus.Labels{
		"name":     "subscription",
		"reporter": "reporter",
	})))

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.chainInfoGauge))
	assert.InDelta(t, 1, testutil.ToFloat64(metricsManager.chainInfoGauge.With(prometheus.Labels{
		"chain":       "chain",
		"pretty_name": "Chain Name",
	})), 0.001)

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.eventsTotalCounter))
	assert.Zero(t, testutil.ToFloat64(metricsManager.eventsTotalCounter.With(prometheus.Labels{
		"chain": "chain",
		"node":  "node",
	})))

	assert.Equal(t, 1, testutil.CollectAndCount(metricsManager.reconnectsCounter))
	assert.Zero(t, testutil.ToFloat64(metricsManager.reconnectsCounter.With(prometheus.Labels{
		"chain": "chain",
		"node":  "node",
	})))
}

func TestMetricsManagerStartDisabled(t *testing.T) {
	t.Parallel()

	config := configPkg.MetricsConfig{Enabled: false}
	logger := loggerPkg.GetNopLogger()
	metricsManager := NewManager(logger, config)
	metricsManager.Start()

	assert.True(t, true)
}

func TestMetricsManagerFailToStart(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	config := configPkg.MetricsConfig{Enabled: true, ListenAddr: "invalid"}
	logger := loggerPkg.GetNopLogger()
	metricsManager := NewManager(logger, config)
	metricsManager.Start()
}

func TestMetricsManagerStopOperation(t *testing.T) {
	t.Parallel()

	config := configPkg.MetricsConfig{Enabled: true, ListenAddr: "invalid"}
	logger := loggerPkg.GetNopLogger()
	metricsManager := NewManager(logger, config)
	metricsManager.Stop()
	assert.True(t, true)
}

//nolint:paralleltest // disabled
func TestAppLoadConfigOk(t *testing.T) {
	config := configPkg.MetricsConfig{Enabled: true, ListenAddr: ":9580"}
	logger := loggerPkg.GetDefaultLogger()
	metricsManager := NewManager(logger, config)
	metricsManager.LogAppVersion("1.2.3")
	go metricsManager.Start()

	for {
		request, err := http.Get("http://localhost:9580/healthcheck")
		_ = request.Body.Close()
		if err == nil {
			break
		}

		time.Sleep(time.Millisecond * 100)
	}

	response, err := http.Get("http://localhost:9580/metrics")
	require.NoError(t, err)
	require.NotEmpty(t, response)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.NotEmpty(t, body)

	err = response.Body.Close()
	require.NoError(t, err)
}
