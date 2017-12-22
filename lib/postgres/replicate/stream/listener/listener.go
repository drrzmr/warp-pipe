package listener

import "context"

// EventListener interface
type EventListener interface {
	Run(ctx context.Context) (err error)
}
