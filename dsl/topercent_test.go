package dsl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToPercent(t *testing.T) {
	require.Equal(t, 0.9, toPercent(90))
	require.Equal(t, 0.99, toPercent(99))
	require.Equal(t, 0.7, toPercent(70))
}
