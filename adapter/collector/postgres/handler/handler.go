package handler

import (
	"unsafe"

	"github.com/jackc/pgx"
	"go.uber.org/zap"

	"github.com/pagarme/warp-pipe/lib/log"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/stream/handler"
	"github.com/pagarme/warp-pipe/pipeline"
)

// Handler interface implementation
type Handler struct {
	handler.EventHandler
	publishCh chan<- pipeline.Message
}

var logger *zap.Logger

func init() { log.Register(&logger, "adapter.collector.postgres.handler") }

// New create new handler
func New(publishCh chan<- pipeline.Message) handler.EventHandler {

	return &Handler{
		publishCh: publishCh,
	}
}

// Heartbeat called for heartbeat event
func (h *Handler) Heartbeat(heartbeat *pgx.ServerHeartbeat) {

	replyRequested := heartbeat.ReplyRequested == 1

	logger.Debug("heartbeat event",
		zap.String("lsn", pgx.FormatLSN(heartbeat.ServerWalEnd)),
		zap.Bool("reply.requested", replyRequested),
	)
}

// Message called for new message
func (h *Handler) Message(msg *pgx.WalMessage) {

	logger.Debug("message event",
		zap.String("lsn", pgx.FormatLSN(msg.WalStart)),
	)

	h.publishCh <- pipeline.NewMessage(pipeline.Payload{
		"WalStart":     msg.WalStart,
		"ServerWalEnd": msg.ServerWalEnd,
		"WalData":      msg.WalData,
		"ServerTime":   msg.ServerTime,
		"ByteLag":      msg.ByteLag(),
	})
}

// Weird called when something strange happen
func (h *Handler) Weird(message *pgx.ReplicationMessage, err error) {

	logger.Panic("weird event",
		zap.Error(err),
		zap.Uintptr("message", uintptr(unsafe.Pointer(message))),
	)
}

// WaitTimeout called for each wait timeout
func (h *Handler) WaitTimeout() { logger.Debug("wait timeout") }

// EOF called when postgres replication connection are closed by peer
func (h *Handler) EOF() { logger.Debug("EOF event"); close(h.publishCh) }
