package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/libtnb/assert/must"

	"github.com/libtnb/fiber-skeleton/internal/shared/registry"
)

func TestCheckReadiness(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		checks := registry.HealthChecks{
			{Name: "one", Check: func(context.Context) error { return nil }},
			{Name: "two", Check: func(context.Context) error { return nil }},
		}
		must.NoError(t, checkReadiness(t.Context(), checks, time.Second))
	})

	t.Run("named failure", func(t *testing.T) {
		checks := registry.HealthChecks{
			{Name: "database", Check: func(context.Context) error { return errors.New("secret backend detail") }},
		}
		err := checkReadiness(t.Context(), checks, time.Second)
		must.ErrorEqual(t, err, "database unavailable")
		must.NotContains(t, err.Error(), "secret")
	})

	t.Run("timeout cancels checker", func(t *testing.T) {
		cancelled := make(chan struct{})
		checks := registry.HealthChecks{
			{Name: "slow", Check: func(ctx context.Context) error {
				<-ctx.Done()
				close(cancelled)
				return context.Cause(ctx)
			}},
		}
		err := checkReadiness(t.Context(), checks, 10*time.Millisecond)
		must.ErrorEqual(t, err, "readiness checks timed out")
		must.Eventually(t, func() bool {
			select {
			case <-cancelled:
				return true
			default:
				return false
			}
		}, must.Tick(time.Millisecond))
	})

	t.Run("failure cancels siblings", func(t *testing.T) {
		cancelled := make(chan struct{})
		checks := registry.HealthChecks{
			{Name: "failed", Check: func(context.Context) error { return errors.New("down") }},
			{Name: "waiting", Check: func(ctx context.Context) error {
				<-ctx.Done()
				close(cancelled)
				return context.Cause(ctx)
			}},
		}
		must.ErrorEqual(t, checkReadiness(t.Context(), checks, time.Second), "failed unavailable")
		must.Eventually(t, func() bool {
			select {
			case <-cancelled:
				return true
			default:
				return false
			}
		}, must.Tick(time.Millisecond))
	})
}
