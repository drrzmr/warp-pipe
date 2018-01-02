package log_test

import (
	"testing"

	"go.uber.org/zap"

	"github.com/pagarme/warp-pipe/lib/log"
)

var logger *zap.Logger

func init() { log.Register(&logger, "test") }

func TestLogDevelopment(t *testing.T) {

	log.Setup(log.Default)
	logger.Debug("test debug message")
}
