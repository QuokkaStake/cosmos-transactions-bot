package toml_config_test

import (
	tomlConfig "main/pkg/config/toml_config"
	"main/pkg/config/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReporterNoName(t *testing.T) {
	t.Parallel()

	reporter := tomlConfig.Reporter{}
	require.Error(t, reporter.Validate())
}

func TestReporterUnsupportedType(t *testing.T) {
	t.Parallel()

	reporter := tomlConfig.Reporter{
		Name: "test",
		Type: "unsupported",
	}
	require.Error(t, reporter.Validate())
}

func TestReporterNoTelegramConfig(t *testing.T) {
	t.Parallel()

	reporter := tomlConfig.Reporter{
		Name: "test",
		Type: "telegram",
	}
	require.Error(t, reporter.Validate())
}

func TestReporterValidTelegram(t *testing.T) {
	t.Parallel()

	reporter := tomlConfig.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &tomlConfig.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
	}
	require.NoError(t, reporter.Validate())
}

func TestReportersInvalid(t *testing.T) {
	t.Parallel()

	reporter := &tomlConfig.Reporter{}
	reporters := tomlConfig.Reporters{reporter}
	require.Error(t, reporters.Validate())
}

func TestReportersDuplicates(t *testing.T) {
	t.Parallel()

	reporter1 := &tomlConfig.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &tomlConfig.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
	}
	reporter2 := &tomlConfig.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &tomlConfig.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
	}
	reporters := tomlConfig.Reporters{reporter1, reporter2}
	require.Error(t, reporters.Validate())
}

func TestReportersValid(t *testing.T) {
	t.Parallel()

	reporter := &tomlConfig.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &tomlConfig.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
	}
	reporters := tomlConfig.Reporters{reporter}
	require.NoError(t, reporters.Validate())
}

func TestHasReporterByName(t *testing.T) {
	t.Parallel()

	reporter := &tomlConfig.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &tomlConfig.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
	}
	reporters := tomlConfig.Reporters{reporter}
	require.True(t, reporters.HasReporterByName("test"))
	require.False(t, reporters.HasReporterByName("test-2"))
}

func TestReporterToAppConfigReporter(t *testing.T) {
	t.Parallel()

	reporter := &tomlConfig.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &tomlConfig.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
	}
	appConfigReporter := reporter.ToAppConfigReporter()

	require.Equal(t, "test", appConfigReporter.Name)
	require.Equal(t, "telegram", appConfigReporter.Type)
	require.Equal(t, int64(1), appConfigReporter.TelegramConfig.Chat)
	require.Equal(t, "xxx:yyy", appConfigReporter.TelegramConfig.Token)
	require.Equal(t, []int64{123}, appConfigReporter.TelegramConfig.Admins)
}

func TestReporterToTomlConfigReporter(t *testing.T) {
	t.Parallel()

	reporter := &types.Reporter{
		Name: "test",
		Type: "telegram",
		TelegramConfig: &types.TelegramConfig{
			Chat:   1,
			Token:  "xxx:yyy",
			Admins: []int64{123},
		},
	}
	tomlConfigReporter := tomlConfig.FromAppConfigReporter(reporter)

	require.Equal(t, "test", tomlConfigReporter.Name)
	require.Equal(t, "telegram", tomlConfigReporter.Type)
	require.Equal(t, int64(1), tomlConfigReporter.TelegramConfig.Chat)
	require.Equal(t, "xxx:yyy", tomlConfigReporter.TelegramConfig.Token)
	require.Equal(t, []int64{123}, tomlConfigReporter.TelegramConfig.Admins)
}
