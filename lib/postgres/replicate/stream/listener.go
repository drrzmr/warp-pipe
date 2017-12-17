package stream

import (
	"context"
	"fmt"

	"github.com/jackc/pgx"
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

func isHeartbeat(m *pgx.ReplicationMessage) bool {
	return m.WalMessage == nil && m.ServerHeartbeat != nil
}
func isMessage(m *pgx.ReplicationMessage) bool          { return m.WalMessage != nil && m.ServerHeartbeat == nil }
func isEOF(err error) bool                              { return err != nil && err.Error() == "EOF" }
func isTimeout(err error) bool                          { return err == context.DeadlineExceeded }
func isCancel(err error) bool                           { return err == context.Canceled }
func isWeird(m *pgx.ReplicationMessage, err error) bool { return (m == nil && err == nil) || err != nil }
