// Package event is the in-process event bus contract: modules publish domain
// events and other modules subscribe, without importing each other. Only the
// interface and its value types live here — the implementation and its provider
// are boot wiring (bootstrap). Swap that implementation for Kafka/NATS later
// and the publishers and subscribers stay unchanged.
package event

import "context"

// Event is a domain event. Name groups subscribers; keep it stable, it is the
// wire identity once a real broker replaces the in-process bus.
type Event interface {
	Name() string
}

// Handler reacts to an event. A handler runs in the publisher's goroutine and
// under its context, so keep it quick or hand off to a job.
type Handler func(ctx context.Context, e Event) error

// Bus publishes events to the handlers subscribed to their name.
type Bus interface {
	Subscribe(name string, h Handler)
	Publish(ctx context.Context, e Event) error
}

// Subscription is the marker a subscriber provider returns; building it
// registers the handler on the bus. A module contributes subscribers under
// registry.SubscriberPrefix, and the app activates them all at startup by
// collecting them — so subscriptions are wired without the publisher and the
// subscriber importing each other.
type Subscription struct{}
