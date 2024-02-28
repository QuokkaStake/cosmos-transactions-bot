package alias_manager

import (
	"os"

	"main/pkg/config"
	configTypes "main/pkg/config/types"

	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog"
)

type Aliases *map[string]string
type TomlAliases map[string]Aliases

type ChainAliases struct {
	Chain   *configTypes.Chain
	Aliases Aliases
}
type AllChainAliases map[string]*ChainAliases

type AliasManager struct {
	Logger  zerolog.Logger
	Path    string
	Chains  configTypes.Chains
	Aliases AllChainAliases
}

func (a AllChainAliases) ToTomlAliases() TomlAliases {
	tomlAliases := make(TomlAliases, len(a))
	for chainName, chainAliases := range a {
		tomlAliases[chainName] = chainAliases.Aliases
	}

	return tomlAliases
}

type ChainAliasesLinks struct {
	Chain *configTypes.Chain
	Links map[string]*configTypes.Link
}

func (a AllChainAliases) ToAliasesLinks() []ChainAliasesLinks {
	aliasesLinks := make([]ChainAliasesLinks, 0)

	for _, chainAliases := range a {
		links := make(map[string]*configTypes.Link)

		if chainAliases.Aliases == nil {
			continue
		}

		for wallet, alias := range *chainAliases.Aliases {
			link := chainAliases.Chain.GetWalletLink(wallet)
			link.Title = alias
			links[wallet] = link
		}

		aliasesLinks = append(aliasesLinks, ChainAliasesLinks{
			Chain: chainAliases.Chain,
			Links: links,
		})
	}

	return aliasesLinks
}

func NewAliasManager(logger *zerolog.Logger, config *config.AppConfig) *AliasManager {
	return &AliasManager{
		Logger:  logger.With().Str("component", "alias_manager").Logger(),
		Path:    config.AliasesPath,
		Chains:  config.Chains,
		Aliases: make(map[string]*ChainAliases, 0),
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

	aliasesBytes, err := os.ReadFile(m.Path)
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

	m.Aliases = make(map[string]*ChainAliases, len(aliasesStruct))
	for chainName, chainAliases := range aliasesStruct {
		chain := m.Chains.FindByName(chainName)
		if chain == nil {
			m.Logger.Fatal().Str("chain", chainName).Msg("Could not find chain found in alias config!")
		}

		m.Aliases[chainName] = &ChainAliases{
			Chain:   chain,
			Aliases: chainAliases,
		}
	}

	m.Logger.Info().Msg("Aliases loaded")
}

func (m *AliasManager) Save() error {
	if !m.Enabled() {
		m.Logger.Warn().Msg("Aliases path not set, not saving aliases")
		return nil
	}

	tomlAliases := m.Aliases.ToTomlAliases()

	f, err := os.Create(m.Path)
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

func (m *AliasManager) Get(chain, address string) string {
	if !m.Enabled() {
		m.Logger.Warn().Msg("Aliases path not set, cannot get alias")
		return ""
	}

	chainAliases, ok := m.Aliases[chain]
	if !ok {
		return ""
	}

	aliases := *chainAliases.Aliases
	alias, ok := aliases[address]
	if !ok {
		return ""
	}

	return alias
}

func (m *AliasManager) Set(chain, address, alias string) error {
	if !m.Enabled() {
		m.Logger.Warn().Msg("Aliases path not set, cannot set alias")
		return nil
	}

	_, ok := m.Aliases[chain]
	if !ok {
		chainFound := m.Chains.FindByName(chain)
		if chainFound == nil {
			m.Logger.Fatal().Str("chain", chain).Msg("Could not find chain when setting an alias!")
		}

		aliases := make(map[string]string, 1)

		m.Aliases[chain] = &ChainAliases{
			Chain:   chainFound,
			Aliases: &aliases,
		}
	}

	chainAliases := m.Aliases[chain]
	aliases := *chainAliases.Aliases
	aliases[address] = alias

	return m.Save()
}

func (m *AliasManager) GetAliasesLinks() []ChainAliasesLinks {
	return m.Aliases.ToAliasesLinks()
}
