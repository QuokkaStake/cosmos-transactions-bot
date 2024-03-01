package reporters

import (
	"main/pkg/alias_manager"
	"main/pkg/config"
	configTypes "main/pkg/config/types"
	"main/pkg/constants"
	"main/pkg/data_fetcher"
	"main/pkg/metrics"
	nodesManager "main/pkg/nodes_manager"
	"main/pkg/reporters/telegram"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type Reporter interface {
	Init()
	Name() string
	Type() string
	Send(report types.Report) error
}

type Reporters []Reporter

func (r Reporters) FindByName(name string) Reporter {
	for _, reporter := range r {
		if reporter.Name() == name {
			return reporter
		}
	}

	return nil
}

func GetReporter(
	reporterConfig *configTypes.Reporter,
	appConfig *config.AppConfig,
	logger *zerolog.Logger,
	nodesManager *nodesManager.NodesManager,
	aliasManager *alias_manager.AliasManager,
	metricsManager *metrics.Manager,
	dataFetcher *data_fetcher.DataFetcher,
	version string,
) Reporter {
	if reporterConfig.Type == constants.ReporterTypeTelegram {
		return telegram.NewReporter(
			reporterConfig,
			appConfig,
			logger,
			nodesManager,
			aliasManager,
			metricsManager,
			dataFetcher,
			version,
		)
	}

	logger.Fatal().Str("type", reporterConfig.Type).Msg("Unsupported reporter received!")
	return nil
}
