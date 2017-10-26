package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSample(t *testing.T) {
	require.Equal(t, "sample", sample())
}
