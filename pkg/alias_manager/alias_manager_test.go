package alias_manager_test

import (
	"errors"
	"main/assets"
	"main/pkg/alias_manager"
	configPkg "main/pkg/config"
	configTypes "main/pkg/config/types"
	"main/pkg/fs"
	loggerPkg "main/pkg/logger"
	"testing"

	"github.com/stretchr/testify/require"
)

type MockFile struct {
	FailWrite bool
	FailClose bool
}

func (file *MockFile) Write(p []byte) (int, error) {
	if file.FailWrite {
		return 1, errors.New("not yet supported")
	}

	return len(p), nil
}

func (file *MockFile) Close() error {
	if file.FailClose {
		return errors.New("not yet supported")
	}

	return nil
}

type MockAliasesFS struct {
	FailCreate bool
	FailWrite  bool
	FailClose  bool
}

func (filesystem *MockAliasesFS) ReadFile(name string) ([]byte, error) {
	return assets.EmbedFS.ReadFile(name)
}

func (filesystem *MockAliasesFS) Create(path string) (fs.File, error) {
	if filesystem.FailCreate {
		return nil, errors.New("not yet supported")
	}

	return &MockFile{
		FailWrite: filesystem.FailWrite,
		FailClose: filesystem.FailClose,
	}, nil
}

func (filesystem *MockAliasesFS) Write(p []byte) (int, error) {
	return 0, errors.New("not yet supported")
}

func TestAliasManagerEnabled(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "path",
	}
	filesystem := &MockAliasesFS{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	require.True(t, aliasManager.Enabled())
}

func TestAliasManagerLoadDisabled(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{}
	filesystem := &MockAliasesFS{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	aliasManager.Load()
	require.Empty(t, aliasManager.Aliases)
}

func TestAliasManagerLoadFailed(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{AliasesPath: "nonexistent.toml"}
	filesystem := &MockAliasesFS{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	aliasManager.Load()
	require.Empty(t, aliasManager.Aliases)
}

func TestAliasManagerLoadInvalidToml(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{AliasesPath: "invalid-toml.toml"}
	filesystem := &MockAliasesFS{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	aliasManager.Load()
	require.Empty(t, aliasManager.Aliases)
}

func TestAliasManagerLoadSuccess(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "valid-aliases.toml",
		Chains: configTypes.Chains{
			{Name: "chain"},
		},
	}
	filesystem := &MockAliasesFS{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	aliasManager.Load()
	require.NotEmpty(t, aliasManager.Aliases)
	require.Equal(t, "alias", aliasManager.Get("subscription", "chain", "wallet"))
}

func TestAliasManagerSaveDisabled(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{}
	filesystem := &MockAliasesFS{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	err := aliasManager.Save()
	require.NoError(t, err)
}

func TestAliasManagerSaveErrorOpening(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{AliasesPath: "savefile.toml"}
	filesystem := &MockAliasesFS{FailCreate: true}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	err := aliasManager.Save()
	require.Error(t, err)
}

func TestAliasManagerSaveErrorWriting(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "savefile.toml",
		Chains: configTypes.Chains{
			{Name: "chain"},
		},
	}
	filesystem := &MockAliasesFS{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	err := aliasManager.Set("subscription", "chain", "wallet", "alias")
	require.NoError(t, err)

	aliasManager.FS = &MockAliasesFS{FailWrite: true}
	err = aliasManager.Save()
	require.Error(t, err)
}

func TestAliasManagerSaveErrorClosing(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "savefile.toml",
		Chains: configTypes.Chains{
			{Name: "chain"},
		},
	}
	filesystem := &MockAliasesFS{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	err := aliasManager.Set("subscription", "chain", "wallet", "alias")
	require.NoError(t, err)

	aliasManager.FS = &MockAliasesFS{FailClose: true}
	err = aliasManager.Save()
	require.Error(t, err)
}

func TestAliasManagerSetDisabled(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{}
	filesystem := &MockAliasesFS{}
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

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "savefile.toml",
		Chains:      configTypes.Chains{},
	}
	filesystem := &MockAliasesFS{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	_ = aliasManager.Set("subscription", "chain", "wallet", "alias")
}

func TestAliasManagerGetLinks(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	config := &configPkg.AppConfig{
		AliasesPath: "savefile.toml",
		Chains: configTypes.Chains{
			{Name: "chain"},
		},
	}
	filesystem := &MockAliasesFS{}
	aliasManager := alias_manager.NewAliasManager(logger, config, filesystem)
	err := aliasManager.Set("subscription", "chain", "wallet", "alias")
	require.NoError(t, err)

	links := aliasManager.GetAliasesLinks("subscription")
	require.Len(t, links, 1)
}
