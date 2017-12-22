package stream

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/stream/handler"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// EventListener interface
type EventListener interface {
	Run(ctx context.Context) (err error)
}

// DefaultEventListener object
type DefaultEventListener struct {
	EventListener
	handler   handler.EventHandler
	replicate *Replicate
}

var listenerLogger = logger.With(zap.String("submodule", "listener"))

// NewDefaultEventListener simple event listener mock
func NewDefaultEventListener(replicate *Replicate, handler handler.EventHandler) EventListener {
	return &DefaultEventListener{
		handler:   handler,
		replicate: replicate,
	}
}

// Run start listener execution
func (d *DefaultEventListener) Run(ctx context.Context) (err error) {

	listenerLogger.Debug("--> Run()")
	defer listenerLogger.Debug("<-- Run()")

	var (
		conn    = d.replicate.conn
		handler = d.handler
		timeout = d.replicate.config.Streaming.WaitMessageTimeout
	)

	for {
		runContext, cancel := context.WithTimeout(ctx, timeout)
		message, err := conn.WaitForReplicationMessage(runContext)
		cancel()

		if ignore, err := filterError(message, handler, err); err != nil {
			return errors.WithStack(err)
		} else if ignore {
			continue
		}

		if isHeartbeat(message) {
			handler.Heartbeat(message.ServerHeartbeat)
			continue
		}

		if isMessage(message) {
			handler.Message(message.WalMessage)
			continue
		}

		handler.Weird(message, err)
	}
}

func filterError(message *pgx.ReplicationMessage, h handler.EventHandler, inErr error) (ignore bool, outErr error) {

	if isTimeout(inErr) {
		h.WaitTimeout()
		return true, nil
	}

	if isCancel(inErr) {
		return false, errors.Wrap(inErr, "canceled context")
	}

	if isEOF(inErr) {
		h.EOF()
		return false, errors.Wrap(inErr, "end of postgres stream messages")
	}

	if isWeird(message, inErr) {
		h.Weird(message, inErr)
		return true, nil
	}

	return false, nil
}

func isHeartbeat(m *pgx.ReplicationMessage) bool {
	return m.WalMessage == nil && m.ServerHeartbeat != nil
}
func isMessage(m *pgx.ReplicationMessage) bool          { return m.WalMessage != nil && m.ServerHeartbeat == nil }
func isEOF(err error) bool                              { return err != nil && err.Error() == "EOF" }
func isTimeout(err error) bool                          { return err == context.DeadlineExceeded }
func isCancel(err error) bool                           { return err == context.Canceled }
func isWeird(m *pgx.ReplicationMessage, err error) bool { return (m == nil && err == nil) || err != nil }
