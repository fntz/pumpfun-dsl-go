package dsl

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculateBuyQuote(t *testing.T) {
	bc := &FetchBoundingCurveResponse{
		RealTokenReserves:    big.NewInt(1000),
		VirtualTokenReserves: big.NewInt(1000),
		VirtualSolReserves:   big.NewInt(1000),
	}

	res := bc.CalculateBuyAmount(100, 0.9)
	require.Equal(t, res, big.NewInt(81))
}
