package log_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/log"
)

func TestLogDevelopment(t *testing.T) {

	logger := log.Development("log")
	require.NotNil(t, logger)
	require.NotNil(t, logger.Logger)

	logger.DebugIf(true, "true")
	logger.DebugIf(false, "false")

	logger.ErrorIf(true, "true")
	logger.ErrorIf(false, "false")

}
