package stream

import (
	"context"
	"fmt"
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
