package yaml_config_test

import (
	"main/pkg/config/types"
	yamlConfig "main/pkg/config/yaml_config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDenomNoName(t *testing.T) {
	t.Parallel()

	denom := yamlConfig.DenomInfo{}
	require.Error(t, denom.Validate())
}

func TestDenomNoDisplayName(t *testing.T) {
	t.Parallel()

	denom := yamlConfig.DenomInfo{
		Denom: "udenom",
	}
	require.Error(t, denom.Validate())
}

func TestDenomValid(t *testing.T) {
	t.Parallel()

	denom := yamlConfig.DenomInfo{
		Denom:        "udenom",
		DisplayDenom: "denom",
	}
	require.NoError(t, denom.Validate())
}

func TestDenomsToAppConfigDenoms(t *testing.T) {
	t.Parallel()

	denom := &yamlConfig.DenomInfo{
		Denom:             "udenom",
		DisplayDenom:      "denom",
		DenomExponent:     10,
		CoingeckoCurrency: "coingecko",
	}
	denoms := yamlConfig.DenomInfos{denom}
	appConfigDenoms := denoms.ToAppConfigDenomInfos()

	require.Len(t, appConfigDenoms, 1)
	require.Equal(t, "udenom", appConfigDenoms[0].Denom)
	require.Equal(t, "denom", appConfigDenoms[0].DisplayDenom)
	require.Equal(t, 10, appConfigDenoms[0].DenomExponent)
	require.Equal(t, "coingecko", appConfigDenoms[0].CoingeckoCurrency)
}

func TestDenomsToYamlConfigDenoms(t *testing.T) {
	t.Parallel()

	denom := &types.DenomInfo{
		Denom:             "udenom",
		DisplayDenom:      "denom",
		DenomExponent:     10,
		CoingeckoCurrency: "coingecko",
	}
	denoms := types.DenomInfos{denom}
	yamlConfigDenoms := yamlConfig.YamlConfigDenomsFrom(denoms)

	require.Len(t, yamlConfigDenoms, 1)
	require.Equal(t, "udenom", yamlConfigDenoms[0].Denom)
	require.Equal(t, "denom", yamlConfigDenoms[0].DisplayDenom)
	require.Equal(t, 10, yamlConfigDenoms[0].DenomExponent)
	require.Equal(t, "coingecko", yamlConfigDenoms[0].CoingeckoCurrency)
}
