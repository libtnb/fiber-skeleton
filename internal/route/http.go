package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/biz"
	"github.com/libtnb/fiber-skeleton/internal/request"
	"github.com/libtnb/fiber-skeleton/internal/service"
)

// HealthRoutes contributes the hello and probe endpoints; none of them are
// documented, so they carry no Request/Response samples.
func HealthRoutes(i do.Injector) (Endpoints, error) {
	health := do.MustInvoke[*service.HealthService](i)

	return Endpoints{
		{Method: fiber.MethodGet, Path: "/", Handler: func(c fiber.Ctx) error {
			return c.SendString("Hello, World 👋!")
		}},
		{Method: fiber.MethodGet, Path: "/healthz", Handler: health.Healthz},
		{Method: fiber.MethodGet, Path: "/readyz", Handler: health.Readyz},
	}, nil
}

// UserRoutes contributes the user endpoints.
func UserRoutes(i do.Injector) (Endpoints, error) {
	user := do.MustInvoke[*service.UserService](i)

	return Endpoints{
		{Method: fiber.MethodGet, Path: "/users", Handler: user.List,
			Summary: "List users", Tags: []string{"user"},
			Request: request.Paginate{}, Response: service.Envelope[service.Page[*biz.User]]{}},
		{Method: fiber.MethodPost, Path: "/users", Handler: user.Create,
			Summary: "Create a user", Tags: []string{"user"},
			Request: request.UserAdd{}, Response: service.Envelope[biz.User]{}},
		{Method: fiber.MethodGet, Path: "/users/:id", Handler: user.Get,
			Summary: "Get a user", Tags: []string{"user"},
			Request: request.UserID{}, Response: service.Envelope[biz.User]{}},
		{Method: fiber.MethodPut, Path: "/users/:id", Handler: user.Update,
			Summary: "Update a user", Tags: []string{"user"},
			Request: request.UserUpdate{}, Response: service.Envelope[biz.User]{}},
		{Method: fiber.MethodDelete, Path: "/users/:id", Handler: user.Delete,
			Summary: "Delete a user", Tags: []string{"user"},
			Request: request.UserID{}},
	}, nil
}
