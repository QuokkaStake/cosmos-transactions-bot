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

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type Manager struct {
	logger zerolog.Logger
	config configPkg.MetricsConfig

	// Chains metrics
	lastBlockHeightCollector   *prometheus.GaugeVec
	lastBlockTimeCollector     *prometheus.GaugeVec
	chainInfoGauge             *prometheus.GaugeVec
	successfulQueriesCollector *prometheus.CounterVec
	failedQueriesCollector     *prometheus.CounterVec
	eventsTotalCounter         *prometheus.CounterVec
	eventsFilteredCounter      *prometheus.CounterVec

	// Node metrics
	nodeConnectedCollector *prometheus.GaugeVec
	reconnectsCounter      *prometheus.CounterVec

	// Reporters metrics
	reporterReportsCounter *prometheus.CounterVec
	reporterErrorsCounter  *prometheus.CounterVec
	reportEntriesCounter   *prometheus.CounterVec
	reporterEnabledGauge   *prometheus.GaugeVec
	reporterQueriesCounter *prometheus.CounterVec

	// Subscriptions metrics
	subscriptionsInfoCounter *prometheus.GaugeVec
	eventsMatchedCounter     *prometheus.CounterVec

	// App metrics
	appVersionGauge *prometheus.GaugeVec
	startTimeGauge  *prometheus.GaugeVec

	registry *prometheus.Registry
}

func NewManager(logger *zerolog.Logger, config configPkg.MetricsConfig) *Manager {
	return &Manager{
		logger: logger.With().Str("component", "metrics").Logger(),
		config: config,

		// Chain metrics
		lastBlockHeightCollector: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "last_height",
			Help: "Height of the last block processed",
		}, []string{"chain"}),
		lastBlockTimeCollector: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "last_time",
			Help: "Time of the last block processed",
		}, []string{"chain"}),
		chainInfoGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "chain_info",
			Help: "Chain info, with constant 1 as value and pretty_name and chain as labels",
		}, []string{"chain", "pretty_name"}),
		successfulQueriesCollector: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "node_successful_queries_total",
			Help: "Counter of successful node queries",
		}, []string{"chain", "node", "type"}),
		failedQueriesCollector: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "node_failed_queries_total",
			Help: "Counter of failed node queries",
		}, []string{"chain", "node", "type"}),
		eventsTotalCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "events_total",
			Help: "WebSocket events received by node",
		}, []string{"chain", "node"}),
		eventsFilteredCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "events_filtered",
			Help: "WebSocket events filtered out by chain, type and reason",
		}, []string{"chain", "type", "reason"}),

		// Node metrics
		nodeConnectedCollector: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "node_connected",
			Help: "Whether the node is successfully connected (1 if yes, 0 if no)",
		}, []string{"chain", "node"}),
		reconnectsCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "reconnects_total",
			Help: "Node reconnects count",
		}, []string{"chain", "node"}),

		// Reporter metrics
		reporterEnabledGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "reporter_enabled",
			Help: "Reporter info, with name and type, always returns 1 as value",
		}, []string{"name", "type"}),
		reporterReportsCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "reporter_reports",
			Help: "Counter of reports sent successfully",
		}, []string{"chain", "reporter", "type", "subscription"}),
		reporterErrorsCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "reporter_errors",
			Help: "Counter of failed reports sends",
		}, []string{"chain", "reporter", "type", "subscription"}),
		reportEntriesCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "report_entries_total",
			Help: "Counter of messages types per each successfully sent report",
		}, []string{"chain", "reporter", "type", "subscription"}),
		reporterQueriesCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "queries",
			Help: "Counter of reporters' queries (like chain status, aliases etc.)",
		}, []string{"reporter", "type"}),

		// Subscription metrics
		subscriptionsInfoCounter: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "subscriptions",
			Help: "Count of chain subscriptions per subscription",
		}, []string{"name", "reporter"}),
		eventsMatchedCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: constants.PrometheusMetricsPrefix + "events_matched",
			Help: "WebSocket events matching filters by chain",
		}, []string{"chain", "type", "subscription"}),

		// App metrics
		appVersionGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "version",
			Help: "App version",
		}, []string{"version"}),
		startTimeGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: constants.PrometheusMetricsPrefix + "start_time",
			Help: "Unix timestamp on when the app was started. Useful for annotations.",
		}, []string{}),
		registry: prometheus.NewRegistry(),
	}
}

func (m *Manager) Start() {
	if !m.config.Enabled {
		m.logger.Info().Msg("Metrics not enabled")
		return
	}

	m.registry.MustRegister(
		m.lastBlockHeightCollector,
		m.lastBlockTimeCollector,
		m.chainInfoGauge,
		m.successfulQueriesCollector,
		m.failedQueriesCollector,
		m.eventsTotalCounter,
		m.eventsFilteredCounter,
		m.nodeConnectedCollector,
		m.reconnectsCounter,
		m.reporterReportsCounter,
		m.reporterErrorsCounter,
		m.reportEntriesCounter,
		m.reporterEnabledGauge,
		m.reporterQueriesCounter,
		m.subscriptionsInfoCounter,
		m.eventsMatchedCounter,
		m.appVersionGauge,
		m.startTimeGauge,
	)

	m.logger.Info().
		Str("addr", m.config.ListenAddr).
		Msg("Metrics handler listening")

	http.Handle("/metrics", promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))
	if err := http.ListenAndServe(m.config.ListenAddr, nil); err != nil {
		m.logger.Fatal().
			Err(err).
			Str("addr", m.config.ListenAddr).
			Msg("Cannot start metrics handler")
	}
}

func (m *Manager) SetAllDefaultMetrics(config *configPkg.AppConfig) {
	m.startTimeGauge.
		With(prometheus.Labels{}).
		Set(float64(time.Now().Unix()))

	for _, chain := range config.Chains {
		m.SetDefaultMetrics(chain)
	}

	for _, subscription := range config.Subscriptions {
		m.subscriptionsInfoCounter.
			With(prometheus.Labels{
				"name":     subscription.Name,
				"reporter": subscription.Reporter,
			}).
			Set(float64(len(subscription.ChainSubscriptions)))
	}
}
func (m *Manager) SetDefaultMetrics(chain *configTypes.Chain) {
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

func (m *Manager) LogQuery(chain string, query queryInfo.QueryInfo, queryType queryInfo.QueryType) {
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

func (m *Manager) LogReport(report types.Report, reporterName string, success bool) {
	if !success {
		m.reporterErrorsCounter.
			With(prometheus.Labels{
				"chain":        report.Chain.Name,
				"reporter":     reporterName,
				"type":         report.Reportable.Type(),
				"subscription": report.Subscription.Name,
			}).
			Inc()
		return
	}

	m.reporterReportsCounter.
		With(prometheus.Labels{
			"chain":        report.Chain.Name,
			"reporter":     reporterName,
			"type":         report.Reportable.Type(),
			"subscription": report.Subscription.Name,
		}).
		Inc()

	for _, entry := range report.Reportable.GetMessages() {
		m.reportEntriesCounter.
			With(prometheus.Labels{
				"chain":        report.Chain.Name,
				"reporter":     reporterName,
				"type":         entry.Type(),
				"subscription": report.Subscription.Name,
			}).
			Inc()
	}
}

func (m *Manager) LogReporterEnabled(name, reporterType string) {
	m.reporterEnabledGauge.
		With(prometheus.Labels{
			"name": name,
			"type": reporterType,
		}).
		Set(1)
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

func (m *Manager) LogFilteredEvent(chain string, eventType string, reason constants.EventFilterReason) {
	m.eventsFilteredCounter.
		With(prometheus.Labels{
			"chain":  chain,
			"type":   eventType,
			"reason": string(reason),
		}).
		Inc()
}

func (m *Manager) LogMatchedEvent(chain string, eventType string, subscription string) {
	m.eventsMatchedCounter.
		With(prometheus.Labels{
			"chain":        chain,
			"type":         eventType,
			"subscription": subscription,
		}).
		Inc()
}

func (m *Manager) LogReporterQuery(reporterName string, query constants.ReporterQuery) {
	m.reporterQueriesCounter.
		With(prometheus.Labels{
			"reporter": reporterName,
			"type":     string(query),
		}).
		Inc()
}

func (m *Manager) LogNodeReconnect(chain string, node string) {
	m.reconnectsCounter.
		With(prometheus.Labels{"chain": chain, "node": node}).
		Inc()
}
