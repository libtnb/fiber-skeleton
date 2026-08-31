// Package registry defines the typed collections assembled by wire.
package registry

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/libtnb/fiber-skeleton/internal/shared/event"
	"github.com/libtnb/fiber-skeleton/internal/shared/job"
	"github.com/libtnb/fiber-skeleton/internal/shared/transport"
)

// Routes contains the endpoint groups contributed by application modules.
type Routes []transport.Endpoints

// Commands contains the management commands contributed by application modules.
type Commands []*cli.Command

// Jobs contains the scheduled jobs contributed by application modules.
type Jobs []job.Fn

// Subscriptions contains the event subscriptions contributed by application modules.
type Subscriptions []event.Subscription

// HealthCheck is one named readiness dependency.
type HealthCheck struct {
	Name  string
	Check func(context.Context) error
}

// HealthChecks contains the readiness dependencies contributed by application modules.
type HealthChecks []HealthCheck
