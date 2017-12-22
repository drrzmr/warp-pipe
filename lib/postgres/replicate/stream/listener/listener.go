package listener

import "context"

// EventListener interface
type EventListener interface {
	Listen(ctx context.Context) (err error)
}
