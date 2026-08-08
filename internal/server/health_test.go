package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
)

func TestCheckReadiness(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		checks := registry.HealthChecks{
			{Name: "one", Check: func(context.Context) error { return nil }},
			{Name: "two", Check: func(context.Context) error { return nil }},
		}
		require.NoError(t, checkReadiness(t.Context(), checks, time.Second))
	})

	t.Run("named failure", func(t *testing.T) {
		checks := registry.HealthChecks{
			{Name: "database", Check: func(context.Context) error { return errors.New("secret backend detail") }},
		}
		err := checkReadiness(t.Context(), checks, time.Second)
		require.EqualError(t, err, "database unavailable")
		require.NotContains(t, err.Error(), "secret")
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
		require.EqualError(t, err, "readiness checks timed out")
		require.Eventually(t, func() bool {
			select {
			case <-cancelled:
				return true
			default:
				return false
			}
		}, time.Second, time.Millisecond)
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
		require.EqualError(t, checkReadiness(t.Context(), checks, time.Second), "failed unavailable")
		require.Eventually(t, func() bool {
			select {
			case <-cancelled:
				return true
			default:
				return false
			}
		}, time.Second, time.Millisecond)
	})
}
