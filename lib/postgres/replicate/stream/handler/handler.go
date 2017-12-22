package handler

import (
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
