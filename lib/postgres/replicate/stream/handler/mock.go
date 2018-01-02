package handler

import (
	"unsafe"

	"github.com/jackc/pgx"
	"go.uber.org/zap"

	"github.com/pagarme/warp-pipe/lib/log"
)

type mockEventHandler struct {
	EventHandler
}

var logger *zap.Logger

func init() { log.Register(&logger, "lib.postgres.replicate.stream.handler") }

// MockEventHandler mock event handler
var MockEventHandler = &mockEventHandler{}

func (m *mockEventHandler) Heartbeat(heartbeat *pgx.ServerHeartbeat) {

	logger.Debug("heartbeat event",
		zap.String("lsn", pgx.FormatLSN(heartbeat.ServerWalEnd)),
		zap.Uint8("reply.requested", heartbeat.ReplyRequested),
	)
}

func (m *mockEventHandler) Message(message *pgx.WalMessage) {

	logger.Debug("message event",
		zap.String("lsn", pgx.FormatLSN(message.WalStart)),
	)
}

func (m *mockEventHandler) WaitTimeout() {

	logger.Debug("wait timeout")
}

func (m *mockEventHandler) EOF() {

	logger.Debug("EOF event")
}

func (m *mockEventHandler) Weird(message *pgx.ReplicationMessage, err error) {

	logger.Panic("weird event",
		zap.Error(err),
		zap.Uintptr("message", uintptr(unsafe.Pointer(message))),
	)
}
