package stream

import (
	"unsafe"

	"github.com/jackc/pgx"
	"go.uber.org/zap"
)

// EventHandler interface
type EventHandler interface {
	Heartbeat(heartbeat *pgx.ServerHeartbeat)
	Message(message *pgx.WalMessage)
	WaitTimeout()
	EOF()
	Weird(message *pgx.ReplicationMessage, err error)
}

type mockEventHandler struct {
	EventHandler
}

var handlerLogger = logger.With(zap.String("submodule", "handler"))

// MockEventHandler mock event handler
var MockEventHandler = &mockEventHandler{}

func (m *mockEventHandler) Heartbeat(heartbeat *pgx.ServerHeartbeat) {

	handlerLogger.Debug("heartbeat event",
		zap.String("lsn", pgx.FormatLSN(heartbeat.ServerWalEnd)),
		zap.Uint8("reply.requested", heartbeat.ReplyRequested),
	)
}

func (m *mockEventHandler) Message(message *pgx.WalMessage) {

	handlerLogger.Debug("message event",
		zap.String("lsn", pgx.FormatLSN(message.WalStart)),
	)
}

func (m *mockEventHandler) WaitTimeout() {

	handlerLogger.Debug("wait timeout")
}

func (m *mockEventHandler) EOF() {

	handlerLogger.Debug("EOF event")
}

func (m *mockEventHandler) Weird(message *pgx.ReplicationMessage, err error) {

	handlerLogger.Panic("weird event",
		zap.Error(err),
		zap.Uintptr("message", uintptr(unsafe.Pointer(message))),
	)
}
