package alias_manager_test

import (
	"main/pkg/alias_manager"
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAliasManagerEnabled(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "path",
	}
	filesystem := &fs.MockFs{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	require.True(t, aliasManager.Enabled())
}

func TestAliasManagerLoadDisabled(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	config := &configPkg.AppConfig{}
	filesystem := &fs.MockFs{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	aliasManager.Load()
	require.Empty(t, aliasManager.Aliases)
}

func TestAliasManagerLoadFailed(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	config := &configPkg.AppConfig{AliasesPath: "nonexistent.toml"}
	filesystem := &fs.MockFs{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	aliasManager.Load()
	require.Empty(t, aliasManager.Aliases)
}

func TestAliasManagerLoadInvalidToml(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	config := &configPkg.AppConfig{AliasesPath: "invalid-toml.toml"}
	filesystem := &fs.MockFs{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	aliasManager.Load()
	require.Empty(t, aliasManager.Aliases)
}

func TestAliasManagerLoadSuccess(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "valid-aliases.toml",
		Chains: configTypes.Chains{
			{Name: "chain"},
		},
	}
	filesystem := &fs.MockFs{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	aliasManager.Load()
	require.NotEmpty(t, aliasManager.Aliases)
	require.Equal(t, "alias", aliasManager.Get("subscription", "chain", "wallet"))
}

func TestAliasManagerSaveDisabled(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	config := &configPkg.AppConfig{}
	filesystem := &fs.MockFs{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	err := aliasManager.Save()
	require.NoError(t, err)
}

func TestAliasManagerSaveErrorOpening(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	config := &configPkg.AppConfig{AliasesPath: "savefile.toml"}
	filesystem := &fs.MockFs{FailCreate: true}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	err := aliasManager.Save()
	require.Error(t, err)
}

func TestAliasManagerSaveErrorWriting(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "savefile.toml",
		Chains: configTypes.Chains{
			{Name: "chain"},
		},
	}
	filesystem := &fs.MockFs{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	err := aliasManager.Set("subscription", "chain", "wallet", "alias")
	require.NoError(t, err)

	aliasManager.FS = &fs.MockFs{FailWrite: true}
	err = aliasManager.Save()
	require.Error(t, err)
}

func TestAliasManagerSaveErrorClosing(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "savefile.toml",
		Chains: configTypes.Chains{
			{Name: "chain"},
		},
	}
	filesystem := &fs.MockFs{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	err := aliasManager.Set("subscription", "chain", "wallet", "alias")
	require.NoError(t, err)

	aliasManager.FS = &fs.MockFs{FailClose: true}
	err = aliasManager.Save()
	require.Error(t, err)
}

func TestAliasManagerSetDisabled(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	config := &configPkg.AppConfig{}
	filesystem := &fs.MockFs{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	err := aliasManager.Set("subscription", "chain", "wallet", "alias")
	require.NoError(t, err)
}

func TestAliasManagerSetNoChain(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	logger := loggerPkg.GetNopLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "savefile.toml",
		Chains:      configTypes.Chains{},
	}
	filesystem := &fs.MockFs{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	_ = aliasManager.Set("subscription", "chain", "wallet", "alias")
}

func TestAliasManagerGetLinks(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "savefile.toml",
		Chains: configTypes.Chains{
			{Name: "chain"},
		},
	}
	filesystem := &fs.MockFs{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	err := aliasManager.Set("subscription", "chain", "wallet", "alias")
	require.NoError(t, err)

	links := aliasManager.GetAliasesLinks("subscription")
	require.Len(t, links, 1)
}
