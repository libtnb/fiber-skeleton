package server

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/libtnb/fiber-skeleton/internal/shared/registry"
	"github.com/libtnb/fiber-skeleton/internal/shared/transport"
)

// HealthRoutes serves the probes; they stay out of the OpenAPI docs.
func HealthRoutes(checks registry.HealthChecks) transport.Endpoints {
	return transport.Endpoints{
		{Method: fiber.MethodGet, Path: "/", Handler: func(c fiber.Ctx) error {
			return c.SendString("Hello, World 👋!")
		}},
		{Method: fiber.MethodGet, Path: "/healthz", Handler: func(c fiber.Ctx) error {
			return c.SendString("ok")
		}},
		{Method: fiber.MethodGet, Path: "/readyz", Handler: func(c fiber.Ctx) error {
			if err := checkReadiness(c.Context(), checks, 5*time.Second); err != nil {
				return transport.Error(c, fiber.StatusServiceUnavailable, "%v", err)
			}
			return c.SendString("ok")
		}},
	}
}

type healthResult struct {
	name string
	err  error
}

func checkReadiness(parent context.Context, checks registry.HealthChecks, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	results := make(chan healthResult, len(checks))
	for _, check := range checks {
		go func() {
			results <- healthResult{name: check.Name, err: check.Check(ctx)}
		}()
	}

	for range checks {
		select {
		case result := <-results:
			if result.err != nil {
				cancel()
				return fmt.Errorf("%s unavailable", result.name)
			}
		case <-ctx.Done():
			return errors.New("readiness checks timed out")
		}
	}
	return nil
}
