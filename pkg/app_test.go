package pkg

import (
	"errors"
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
	filtererPkg "main/pkg/filterer"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"
	metricsPkg "main/pkg/metrics"
	reportersPkg "main/pkg/reporters"
	"main/pkg/types"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestNewAppInvalidYaml(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	NewApp(&fs.MockFs{}, "invalid-yaml.yml", "1.2.3")
}

func TestNewAppValidYamlWithWarnings(t *testing.T) {
	t.Parallel()

	app := NewApp(&fs.MockFs{}, "valid-unused-chain.yml", "1.2.3")
	require.NotNil(t, app)
}

func TestNewAppValidYaml(t *testing.T) {
	t.Parallel()

	app := NewApp(&fs.MockFs{}, "valid.yml", "1.2.3")
	require.NotNil(t, app)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestAppStart(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "http://localhost:9580/healthcheck", httpmock.InitialTransport.RoundTrip)
	httpmock.RegisterResponder("GET", "http://localhost:9580/metrics", httpmock.InitialTransport.RoundTrip)

	app := NewApp(&fs.MockFs{}, "valid.yml", "1.2.3")
	require.NotNil(t, app)

	app.Reporters = reportersPkg.Reporters{
		&reportersPkg.TestReporter{ReporterName: "test-reporter"},
		&reportersPkg.TestReporter{ReporterName: "test-reporter-2", FailToSend: true},
		&reportersPkg.TestReporter{ReporterName: "test-reporter-3", FailToInit: true},
	}

	go app.Start()

	for {
		request, err := http.Get("http://localhost:9580/healthcheck")
		if request != nil && request.Body != nil {
			_ = request.Body.Close()
		}
		if err == nil {
			break
		}

		time.Sleep(time.Millisecond * 100)
	}

	app.NodesManager.Channel <- types.Report{
		Chain: &configTypes.Chain{Name: "chain"},
		Node:  "node",
		Reportable: &types.NodeConnectError{
			Error: errors.New("some error"),
		},
	}

	app.QuitChannel <- syscall.SIGTERM
}

func TestAppProcessReport(t *testing.T) {
	t.Parallel()

	config := &configPkg.AppConfig{
		Chains: configTypes.Chains{
			{Name: "chain"},
		},
		Subscriptions: configTypes.Subscriptions{
			{
				Name:     "subscription",
				Reporter: "test-reporter",
				ChainSubscriptions: configTypes.ChainSubscriptions{
					{
						Chain:         "chain",
						Filters:       configTypes.Filters{},
						LogNodeErrors: true,
					},
				},
			},
			{
				Name:     "subscription",
				Reporter: "test-reporter-2",
				ChainSubscriptions: configTypes.ChainSubscriptions{
					{
						Chain:         "chain",
						Filters:       configTypes.Filters{},
						LogNodeErrors: true,
					},
				},
			},
		},
		Reporters: configTypes.Reporters{
			{
				Name: "test-reporter",
			},
			{
				Name: "test-reporter-2",
			},
		},
	}

	logger := loggerPkg.GetNopLogger()
	metricsManager := metricsPkg.NewManager(logger, configPkg.MetricsConfig{Enabled: true})
	filterer := filtererPkg.NewFilterer(logger, config, metricsManager)

	app := &App{
		Reporters: reportersPkg.Reporters{
			&reportersPkg.TestReporter{ReporterName: "test-reporter"},
			&reportersPkg.TestReporter{ReporterName: "test-reporter-2", FailToSend: true},
			&reportersPkg.TestReporter{ReporterName: "test-reporter-3", FailToInit: true},
		},
		Filterer:       filterer,
		MetricsManager: metricsManager,
	}

	report := types.Report{
		Chain: &configTypes.Chain{Name: "chain"},
		Node:  "node",
		Reportable: &types.NodeConnectError{
			Error: errors.New("some error"),
		},
	}

	app.ProcessReport(report)
}
