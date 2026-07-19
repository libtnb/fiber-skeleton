// Package event is the bus contract: modules publish and subscribe without
// importing each other; the implementation lives in bootstrap.
package event

import "context"

// Event is a domain event; Name is its wire identity, keep it stable.
type Event interface {
	Name() string
}

// Handler runs in the publisher's goroutine; keep it quick.
type Handler func(ctx context.Context, e Event) error

// Bus publishes events to the handlers subscribed to their name.
type Bus interface {
	Subscribe(name string, h Handler)
	Publish(ctx context.Context, e Event) error
}

// Subscription is the marker a subscriber provider returns; building it
// registers the handler on the bus.
type Subscription struct{}
