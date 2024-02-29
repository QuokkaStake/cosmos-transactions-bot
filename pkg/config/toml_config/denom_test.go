package toml_config_test

import (
	tomlConfig "main/pkg/config/toml_config"
	"main/pkg/config/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDenomNoName(t *testing.T) {
	t.Parallel()

	denom := tomlConfig.DenomInfo{}
	require.Error(t, denom.Validate())
}

func TestDenomNoDisplayName(t *testing.T) {
	t.Parallel()

	denom := tomlConfig.DenomInfo{
		Denom: "udenom",
	}
	require.Error(t, denom.Validate())
}

func TestDenomValid(t *testing.T) {
	t.Parallel()

	denom := tomlConfig.DenomInfo{
		Denom:        "udenom",
		DisplayDenom: "denom",
	}
	require.NoError(t, denom.Validate())
}

func TestDenomsToAppConfigDenoms(t *testing.T) {
	t.Parallel()

	denom := &tomlConfig.DenomInfo{
		Denom:             "udenom",
		DisplayDenom:      "denom",
		DenomExponent:     10,
		CoingeckoCurrency: "coingecko",
	}
	denoms := tomlConfig.DenomInfos{denom}
	appConfigDenoms := denoms.ToAppConfigDenomInfos()

	require.Len(t, appConfigDenoms, 1)
	require.Equal(t, "udenom", appConfigDenoms[0].Denom)
	require.Equal(t, "denom", appConfigDenoms[0].DisplayDenom)
	require.Equal(t, 10, appConfigDenoms[0].DenomExponent)
	require.Equal(t, "coingecko", appConfigDenoms[0].CoingeckoCurrency)
}

func TestDenomsToTomlConfigDenoms(t *testing.T) {
	t.Parallel()

	denom := &types.DenomInfo{
		Denom:             "udenom",
		DisplayDenom:      "denom",
		DenomExponent:     10,
		CoingeckoCurrency: "coingecko",
	}
	denoms := types.DenomInfos{denom}
	tomlConfigDenoms := tomlConfig.TomlConfigDenomsFrom(denoms)

	require.Len(t, tomlConfigDenoms, 1)
	require.Equal(t, "udenom", tomlConfigDenoms[0].Denom)
	require.Equal(t, "denom", tomlConfigDenoms[0].DisplayDenom)
	require.Equal(t, 10, tomlConfigDenoms[0].DenomExponent)
	require.Equal(t, "coingecko", tomlConfigDenoms[0].CoingeckoCurrency)
}
