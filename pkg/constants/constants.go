package constants

type EventFilterReason string

type ReporterQuery string

const (
	PrometheusMetricsPrefix string = "cosmos_transactions_bot_"

	ReporterTypeTelegram string = "telegram"

	EventFilterReasonTxErrorNotLogged            EventFilterReason = "tx_error_not_logged"
	EventFilterReasonNodeErrorNotLogged          EventFilterReason = "node_error_not_logged"
	EventFilterReasonUnsupportedMsgTypeNotLogged EventFilterReason = "unsupported_msg_type_not_logged"
	EventFilterReasonFailedTxNotLogged           EventFilterReason = "failed_tx_not_logged"
	EventFilterReasonEmptyTxNotLogged            EventFilterReason = "empty_tx_not_logged"

	ReporterQueryHelp        ReporterQuery = "help"
	ReporterQueryGetAliases  ReporterQuery = "get_aliases"
	ReporterQuerySetAlias    ReporterQuery = "set_alias"
	ReporterQueryNodesStatus ReporterQuery = "nodes_status"
)

func GetReporterTypes() []string {
	return []string{
		ReporterTypeTelegram,
	}
}
