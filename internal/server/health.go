package server

import (
	"github.com/gofiber/fiber/v3"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
)

// HealthRoutes serves the probes; they stay out of the OpenAPI docs.
func HealthRoutes(i do.Injector) (transport.Endpoints, error) {
	return transport.Endpoints{
		{Method: fiber.MethodGet, Path: "/", Handler: func(c fiber.Ctx) error {
			return c.SendString("Hello, World 👋!")
		}},
		{Method: fiber.MethodGet, Path: "/healthz", Handler: func(c fiber.Ctx) error {
			return c.SendString("ok")
		}},
		{Method: fiber.MethodGet, Path: "/readyz", Handler: func(c fiber.Ctx) error {
			for name, err := range i.HealthCheckWithContext(c.Context()) {
				if err != nil {
					return transport.Error(c, fiber.StatusServiceUnavailable, "%s unavailable", name)
				}
			}
			return c.SendString("ok")
		}},
	}, nil
}
