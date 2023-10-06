package metrics

import (
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
	"main/pkg/constants"
	"main/pkg/types"
	queryInfo "main/pkg/types/query_info"
	"main/pkg/utils"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Manager struct {
	logger zerolog.Logger
	config configPkg.MetricsConfig

	lastBlockHeightCollector *prometheus.GaugeVec
	lastBlockTimeCollector   *prometheus.GaugeVec
	nodeConnectedCollector   *prometheus.GaugeVec
	reconnectsCounter        *prometheus.CounterVec

	successfulQueriesCollector *prometheus.CounterVec
	failedQueriesCollector     *prometheus.CounterVec

	eventsTotalCounter    *prometheus.CounterVec
	eventsFilteredCounter *prometheus.CounterVec
	eventsMatchedCounter  *prometheus.CounterVec

	reportsCounter       *prometheus.CounterVec
	reportEntriesCounter *prometheus.CounterVec

	reporterEnabledGauge   *prometheus.GaugeVec
	reporterQueriesCounter *prometheus.CounterVec

	appVersionGauge *prometheus.GaugeVec
	chainInfoGauge  *prometheus.GaugeVec
}

func NewManager(logger *zerolog.Logger, config configPkg.MetricsConfig) *Manager {
	return &Manager{
		logger: logger.With().Str("component", "metrics").Logger(),
		config: config,
		lastBlockHeightCollector: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "last_height",
			Help: "Height of the last block processed",
		}, []string{"chain"}),
		lastBlockTimeCollector: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "last_time",
			Help: "Time of the last block processed",
		}, []string{"chain"}),
		nodeConnectedCollector: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "node_connected",
			Help: "Whether the node is successfully connected (1 if yes, 0 if no)",
		}, []string{"chain", "node"}),
		successfulQueriesCollector: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "node_successful_queries_total",
			Help: "Counter of successful node queries",
		}, []string{"chain", "node", "type"}),
		failedQueriesCollector: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "node_failed_queries_total",
			Help: "Counter of failed node queries",
		}, []string{"chain", "node", "type"}),
		reportsCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "node_reports",
			Help: "Counter of reports send",
		}, []string{"chain"}),
		reportEntriesCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "node_report_entries_total",
			Help: "Counter of report entries send",
		}, []string{"chain", "type"}),
		reporterEnabledGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "reporter_enabled",
			Help: "Whether the reporter is enabled (1 if yes, 0 if no)",
		}, []string{"name"}),
		reporterQueriesCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "reporter_queries",
			Help: "Reporters' queries count ",
		}, []string{"chain", "name", "query"}),
		appVersionGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "version",
			Help: "App version",
		}, []string{"version"}),
		eventsTotalCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "events_total",
			Help: "WebSocket events received by node",
		}, []string{"chain", "node"}),
		eventsFilteredCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "events_filtered",
			Help: "WebSocket events filtered out by chain",
		}, []string{"chain", "type"}),
		eventsMatchedCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "events_matched",
			Help: "WebSocket events matching filters by chain",
		}, []string{"chain", "type"}),
		reconnectsCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "reconnects_total",
			Help: "Node reconnects count",
		}, []string{"chain", "node"}),
		chainInfoGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "chain_info",
			Help: "Chain info, with constant 1 as value and pretty_name and chain as labels",
		}, []string{"chain", "pretty_name"}),
	}
}

func (m *Manager) SetAllDefaultMetrics(chains []*configTypes.Chain) {
	for _, chain := range chains {
		m.SetDefaultMetrics(chain)
	}
}
func (m *Manager) SetDefaultMetrics(chain *configTypes.Chain) {
	m.reportsCounter.
		With(prometheus.Labels{"chain": chain.Name}).
		Add(0)

	m.chainInfoGauge.
		With(prometheus.Labels{"chain": chain.Name, "pretty_name": chain.PrettyName}).
		Set(1)

	for _, node := range chain.TendermintNodes {
		m.eventsTotalCounter.
			With(prometheus.Labels{"chain": chain.Name, "node": node}).
			Add(0)

		m.reconnectsCounter.
			With(prometheus.Labels{"chain": chain.Name, "node": node}).
			Add(0)
	}

	for _, node := range chain.TendermintNodes {
		m.eventsTotalCounter.
			With(prometheus.Labels{"chain": chain.Name, "node": node}).
			Add(0)

		m.reconnectsCounter.
			With(prometheus.Labels{"chain": chain.Name, "node": node}).
			Add(0)

		for _, queryType := range queryInfo.GetQueryTypes() {
			m.successfulQueriesCollector.
				With(prometheus.Labels{"chain": chain.Name, "node": node, "type": string(queryType)}).
				Add(0)

			m.failedQueriesCollector.
				With(prometheus.Labels{"chain": chain.Name, "node": node, "type": string(queryType)}).
				Add(0)
		}
	}
}

func (m *Manager) Start() {
	if !m.config.Enabled {
		m.logger.Info().Msg("Metrics not enabled")
		return
	}

	m.logger.Info().
		Str("addr", m.config.ListenAddr).
		Msg("Metrics handler listening")

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(m.config.ListenAddr, nil); err != nil {
		m.logger.Fatal().
			Err(err).
			Str("addr", m.config.ListenAddr).
			Msg("Cannot start metrics handler")
	}
}

func (m *Manager) LogLastHeight(chain string, height int64, blockTime time.Time) {
	m.lastBlockHeightCollector.
		With(prometheus.Labels{"chain": chain}).
		Set(float64(height))

	m.lastBlockTimeCollector.
		With(prometheus.Labels{"chain": chain}).
		Set(float64(blockTime.Unix()))
}

func (m *Manager) LogNodeConnection(chain, node string, connected bool) {
	m.nodeConnectedCollector.
		With(prometheus.Labels{"chain": chain, "node": node}).
		Set(utils.BoolToFloat64(connected))
}

func (m *Manager) LogTendermintQuery(chain string, query queryInfo.QueryInfo, queryType queryInfo.QueryType) {
	if query.Success {
		m.successfulQueriesCollector.
			With(prometheus.Labels{
				"chain": chain,
				"node":  query.Node,
				"type":  string(queryType),
			}).Inc()
	} else {
		m.failedQueriesCollector.
			With(prometheus.Labels{
				"chain": chain,
				"node":  query.Node,
				"type":  string(queryType),
			}).Inc()
	}
}

func (m *Manager) LogReport(report types.Report) {
	m.reportsCounter.
		With(prometheus.Labels{"chain": report.Chain.Name}).
		Inc()

	for _, entry := range report.Reportable.GetMessages() {
		m.reportEntriesCounter.
			With(prometheus.Labels{
				"chain": report.Chain.Name,
				"type":  entry.Type(),
			}).
			Inc()
	}
}

func (m *Manager) LogReporterEnabled(name string, enabled bool) {
	m.reporterEnabledGauge.
		With(prometheus.Labels{"name": name}).
		Set(utils.BoolToFloat64(enabled))
}

func (m *Manager) LogAppVersion(version string) {
	m.appVersionGauge.
		With(prometheus.Labels{"version": version}).
		Set(1)
}

func (m *Manager) LogWSEvent(chain string, node string) {
	m.eventsTotalCounter.
		With(prometheus.Labels{"chain": chain, "node": node}).
		Inc()
}

func (m *Manager) LogFilteredEvent(chain string, eventType string) {
	m.eventsFilteredCounter.
		With(prometheus.Labels{"chain": chain, "type": eventType}).
		Inc()
}

func (m *Manager) LogMatchedEvent(chain string, eventType string) {
	m.eventsMatchedCounter.
		With(prometheus.Labels{"chain": chain, "type": eventType}).
		Inc()
}

func (m *Manager) LogNodeReconnect(chain string, node string) {
	m.reconnectsCounter.
		With(prometheus.Labels{"chain": chain, "node": node}).
		Inc()
}
