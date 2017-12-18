package stream

import (
	"fmt"

	"github.com/jackc/pgx"
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

// MockEventHandler mock event handler
var MockEventHandler = &mockEventHandler{}

func (m *mockEventHandler) Heartbeat(heartbeat *pgx.ServerHeartbeat) {
	fmt.Println("[handler] heartbeat")
}

func (m *mockEventHandler) Message(message *pgx.WalMessage) {

	fmt.Printf("[handler] message, %#x\n", message.WalStart)
}

func (m *mockEventHandler) WaitTimeout() {
	fmt.Println("[handler] wait timeout")
}

func (m *mockEventHandler) EOF() {
	fmt.Println("[handler] EOF")
}

func (m *mockEventHandler) Weird(message *pgx.ReplicationMessage, err error) {
	fmt.Println("[handler] weird, message:", message, ", error:", err)
}
