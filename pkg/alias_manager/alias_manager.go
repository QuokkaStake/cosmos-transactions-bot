package alias_manager

import (
	"main/pkg/config"
	configTypes "main/pkg/config/types"
	"main/pkg/fs"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
)

type AliasManager struct {
	Logger  zerolog.Logger
	Path    string
	Chains  configTypes.Chains
	Aliases AllAliases
	FS      fs.FS
}

func NewAliasManager(
	logger *zerolog.Logger,
	config *config.AppConfig,
	fs fs.FS,
) *AliasManager {
	return &AliasManager{
		Logger:  logger.With().Str("component", "alias_manager").Logger(),
		Path:    config.AliasesPath,
		Chains:  config.Chains,
		Aliases: AllAliases{},
		FS:      fs,
	}
}

func (m *AliasManager) Enabled() bool {
	return m.Path != ""
}

func (m *AliasManager) Load() {
	if !m.Enabled() {
		m.Logger.Warn().Msg("Aliases path not set, not loading aliases")
		return
	}

	aliasesBytes, err := m.FS.ReadFile(m.Path)
	if err != nil {
		m.Logger.Error().Err(err).Msg("Could not load aliases")
		return
	}

	aliasesString := string(aliasesBytes)

	var aliasesStruct TomlAliases
	if _, err = toml.Decode(aliasesString, &aliasesStruct); err != nil {
		m.Logger.Error().Err(err).Msg("Could not decode aliases")
		return
	}

	m.Aliases = aliasesStruct.ToAliases(m.Chains, m.Logger)
	m.Logger.Info().Msg("Aliases loaded")
}

func (m *AliasManager) Save() error {
	if !m.Enabled() {
		m.Logger.Warn().Msg("Aliases path not set, not saving aliases")
		return nil
	}

	tomlAliases := m.Aliases.ToTomlAliases()

	f, err := m.FS.Create(m.Path)
	if err != nil {
		m.Logger.Error().Err(err).Msg("Could not create aliases file")
		return err
	}
	if err := toml.NewEncoder(f).Encode(tomlAliases); err != nil {
		m.Logger.Error().Err(err).Msg("Could not save aliases")
		return err
	}
	if err := f.Close(); err != nil {
		m.Logger.Error().Err(err).Msg("Could not close aliases file when saving")
		return err
	}

	return nil
}

func (m *AliasManager) Get(subscription, chain, address string) string {
	return m.Aliases.Get(subscription, chain, address)
}

func (m *AliasManager) Set(subscription, chainName, address, alias string) error {
	if !m.Enabled() {
		m.Logger.Warn().Msg("Aliases path not set, cannot set alias")
		return nil
	}

	chainFound := m.Chains.FindByName(chainName)
	if chainFound == nil {
		m.Logger.Panic().
			Str("chain", chainName).
			Msg("Could not find chain when setting an alias!")
	}

	m.Aliases.Set(subscription, chainFound, address, alias)
	return m.Save()
}

func (m *AliasManager) GetAliasesLinks(subscription string) []ChainAliasesLinks {
	return m.Aliases.GetAliasesLinks(subscription)
}
