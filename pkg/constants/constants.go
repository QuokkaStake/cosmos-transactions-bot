package constants

const (
	PrometheusMetricsPrefix = "cosmos_transactions_bot_"

	ReporterTypeTelegram string = "telegram"
)

func GetReporterTypes() []string {
	return []string{
		ReporterTypeTelegram,
	}
}
