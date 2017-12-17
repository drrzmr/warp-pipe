package stream

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

// EventListener interface
type EventListener interface {
	Run(ctx context.Context) (err error)
}

// DefaultEventListener object
type DefaultEventListener struct {
	EventListener
	handler   EventHandler
	replicate *Replicate

	consumedWalPosition uint64
}

// NewDefaultEventListener simple event listener mock
func NewDefaultEventListener(replicate *Replicate, handler EventHandler) EventListener {
	return &DefaultEventListener{
		handler:   handler,
		replicate: replicate,
	}
}

// Run start listener execution
func (d *DefaultEventListener) Run(ctx context.Context) (err error) {

	fmt.Printf("[listener] Run() -->\n")

	d.handler.EOF()
	if ctx != nil {
		fmt.Println(ctx.Value("name"))
	}

	fmt.Printf("[listener] Run() <--\n")
	return nil
}

func filterError(message *pgx.ReplicationMessage, handler EventHandler, inErr error) (ignore bool, outErr error) {

	if isTimeout(inErr) {
		handler.WaitTimeout()
		return true, nil
	}

	if isCancel(inErr) {
		return false, errors.Wrap(inErr, "canceled context")
	}

	if isEOF(inErr) {
		handler.EOF()
		return false, errors.Wrap(inErr, "end of postgres stream messages")
	}

	if isWeird(message, inErr) {
		handler.Weird(message, inErr)
		return true, nil
	}

	return false, nil
}

func sendStandByStatus(conn *pgx.ReplicationConn, consumedWalPosition *uint64) (err error) {

	var (
		status   *pgx.StandbyStatus
		position = atomic.LoadUint64(consumedWalPosition)
	)

	if status, err = pgx.NewStandbyStatus(position); err != nil {
		return errors.Wrapf(err, "create new standby status object failed, position: %d", position)
	}

	err = conn.SendStandbyStatus(status)
	if err == nil {
		fmt.Printf("[listener] send standby status, position: %d\n", position)
	}
	return errors.Wrapf(err, "send stand by status failed, position: %d", position)
}

func isHeartbeat(m *pgx.ReplicationMessage) bool {
	return m.WalMessage == nil && m.ServerHeartbeat != nil
}
func isMessage(m *pgx.ReplicationMessage) bool          { return m.WalMessage != nil && m.ServerHeartbeat == nil }
func isEOF(err error) bool                              { return err != nil && err.Error() == "EOF" }
func isTimeout(err error) bool                          { return err == context.DeadlineExceeded }
func isCancel(err error) bool                           { return err == context.Canceled }
func isWeird(m *pgx.ReplicationMessage, err error) bool { return (m == nil && err == nil) || err != nil }