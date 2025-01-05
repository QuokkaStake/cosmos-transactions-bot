package alias_manager

import (
	"main/pkg/config"
	configTypes "main/pkg/config/types"
	"main/pkg/fs"

	"gopkg.in/yaml.v3"

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

	var aliasesStruct YamlAliases
	if err = yaml.Unmarshal(aliasesBytes, &aliasesStruct); err != nil {
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

	yamlAliases := m.Aliases.ToYamlAliases()

	f, err := m.FS.Create(m.Path)
	if err != nil {
		m.Logger.Error().Err(err).Msg("Could not create aliases file")
		return err
	}
	if encodeErr := yaml.NewEncoder(f).Encode(yamlAliases); encodeErr != nil {
		m.Logger.Error().Err(encodeErr).Msg("Could not save aliases")
		return encodeErr
	}
	if closeErr := f.Close(); closeErr != nil {
		m.Logger.Error().Err(closeErr).Msg("Could not close aliases file when saving")
		return closeErr
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
