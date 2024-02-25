package types_test

import (
	"main/pkg/config/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDenomsFind(t *testing.T) {
	t.Parallel()

	denoms := types.DenomInfos{
		{Denom: "denom"},
	}

	require.NotNil(t, denoms.Find("denom"))
	require.Nil(t, denoms.Find("denom-2"))
}
