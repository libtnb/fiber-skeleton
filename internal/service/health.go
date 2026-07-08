package service

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/samber/do/v2"
)

// healthChecker is the one container capability Readyz needs; keeping the
// dependency this narrow makes the service easy to fake in tests.
type healthChecker interface {
	HealthCheckWithContext(ctx context.Context) map[string]error
}

// HealthService serves the container/orchestrator probes.
type HealthService struct {
	checker healthChecker
}

func NewHealthService(i do.Injector) (*HealthService, error) {
	return &HealthService{
		checker: i,
	}, nil
}

// Healthz is the liveness probe: the process is up and serving.
func (r *HealthService) Healthz(c fiber.Ctx) error {
	return c.SendString("ok")
}

// Readyz is the readiness probe: every health-checkable service in the
// container (the database, and whatever you add later) must pass.
func (r *HealthService) Readyz(c fiber.Ctx) error {
	for name, err := range r.checker.HealthCheckWithContext(c.Context()) {
		if err != nil {
			return Error(c, fiber.StatusServiceUnavailable, "%s unavailable", name)
		}
	}

	return c.SendString("ok")
}
