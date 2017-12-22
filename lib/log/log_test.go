package log_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/log"
)

func TestLogDevelopment(t *testing.T) {

	logger := log.Development("log")
	require.NotNil(t, logger)
}
