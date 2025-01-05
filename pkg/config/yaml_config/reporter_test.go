package yaml_config_test

import (
	"main/pkg/config/types"
	yamlConfig "main/pkg/config/yaml_config"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReporterNoName(t *testing.T) {
	t.Parallel()

	reporter := yamlConfig.Reporter{}
	require.Error(t, reporter.Validate())
}

func TestYamlConfigInvalidTimezone(t *testing.T) {
	t.Parallel()

	reporter := yamlConfig.Reporter{Name: "test", Timezone: "invalid"}
	require.Error(t, reporter.Validate())
}

func TestReporterUnsupportedType(t *testing.T) {
	t.Parallel()

	reporter := yamlConfig.Reporter{
		Name:     "test",
		Type:     "unsupported",
		Timezone: "Etc/GMT",
	}
	require.Error(t, reporter.Validate())
}

func TestReporterNoTelegramConfig(t *testing.T) {
	t.Parallel()

	reporter := yamlConfig.Reporter{
		Name:     "test",
		Type:     "telegram",
		Timezone: "Etc/GMT",
	}
	require.Error(t, reporter.Validate())
}

func TestReporterValidTelegram(t *testing.T) {
	t.Parallel()

	reporter := yamlConfig.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &yamlConfig.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
		Timezone: "Etc/GMT",
	}
	require.NoError(t, reporter.Validate())
}

func TestReportersInvalid(t *testing.T) {
	t.Parallel()

	reporter := &yamlConfig.Reporter{}
	reporters := yamlConfig.Reporters{reporter}
	require.Error(t, reporters.Validate())
}

func TestReportersDuplicates(t *testing.T) {
	t.Parallel()

	reporter1 := &yamlConfig.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &yamlConfig.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
	}
	reporter2 := &yamlConfig.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &yamlConfig.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
	}
	reporters := yamlConfig.Reporters{reporter1, reporter2}
	require.Error(t, reporters.Validate())
}

func TestReportersValid(t *testing.T) {
	t.Parallel()

	reporter := &yamlConfig.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &yamlConfig.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
	}
	reporters := yamlConfig.Reporters{reporter}
	require.NoError(t, reporters.Validate())
}

func TestHasReporterByName(t *testing.T) {
	t.Parallel()

	reporter := &yamlConfig.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &yamlConfig.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
	}
	reporters := yamlConfig.Reporters{reporter}
	require.True(t, reporters.HasReporterByName("test"))
	require.False(t, reporters.HasReporterByName("test-2"))
}

func TestReporterToAppConfigReporter(t *testing.T) {
	t.Parallel()

	reporter := &yamlConfig.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &yamlConfig.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
		Timezone: "Etc/GMT",
	}
	appConfigReporter := reporter.ToAppConfigReporter()

	require.Equal(t, "test", appConfigReporter.Name)
	require.Equal(t, "telegram", appConfigReporter.Type)
	require.Equal(t, int64(1), appConfigReporter.TelegramConfig.Chat)
	require.Equal(t, "xxx:yyy", appConfigReporter.TelegramConfig.Token)
	require.Equal(t, []int64{123}, appConfigReporter.TelegramConfig.Admins)
	require.Equal(t, "Etc/GMT", appConfigReporter.Timezone.String())
}

func TestReporterToYamlConfigReporter(t *testing.T) {
	t.Parallel()

	timezone, err := time.LoadLocation("Etc/GMT")
	require.NoError(t, err)

	reporter := &types.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &types.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
		Timezone: timezone,
	}
	yamlConfigReporter := yamlConfig.FromAppConfigReporter(reporter)

	require.Equal(t, "test", yamlConfigReporter.Name)
	require.Equal(t, "telegram", yamlConfigReporter.Type)
	require.Equal(t, int64(1), yamlConfigReporter.TelegramConfig.Chat)
	require.Equal(t, "xxx:yyy", yamlConfigReporter.TelegramConfig.Token)
	require.Equal(t, []int64{123}, yamlConfigReporter.TelegramConfig.Admins)
	require.Equal(t, "Etc/GMT", yamlConfigReporter.Timezone)
}
