package listener

import (
	"context"
	"time"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"

	"github.com/pagarme/warp-pipe/lib/log"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/stream/handler"
)

// DefaultEventListener object
type DefaultEventListener struct {
	EventListener
	handler handler.EventHandler
	conn    *pgx.ReplicationConn
	timeout time.Duration
}

var logger = log.Development("listener")

// NewDefaultEventListener simple event listener mock
func NewDefaultEventListener(conn *pgx.ReplicationConn, timeout time.Duration, h handler.EventHandler) EventListener {
	return &DefaultEventListener{
		handler: h,
		conn:    conn,
		timeout: timeout,
	}
}

// Listen start listener execution
func (d *DefaultEventListener) Listen(ctx context.Context) (err error) {

	logger.Debug("--> Run()")
	defer logger.Debug("<-- Run()")

	for {
		runContext, cancel := context.WithTimeout(ctx, d.timeout)
		message, err := d.conn.WaitForReplicationMessage(runContext)
		cancel()

		if ignore, err := filterError(message, d.handler, err); err != nil {
			return errors.WithStack(err)
		} else if ignore {
			continue
		}

		if isHeartbeat(message) {
			d.handler.Heartbeat(message.ServerHeartbeat)
			continue
		}

		if isMessage(message) {
			d.handler.Message(message.WalMessage)
			continue
		}

		d.handler.Weird(message, err)
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
