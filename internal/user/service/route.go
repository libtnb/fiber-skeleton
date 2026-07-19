package service

import (
	"github.com/gofiber/fiber/v3"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
	"github.com/libtnb/fiber-skeleton/internal/user/biz"
)

func UserRoutes(i do.Injector) (transport.Endpoints, error) {
	user := do.MustInvoke[*UserService](i)

	return transport.Endpoints{
		{Method: fiber.MethodGet, Path: "/users", Handler: user.List,
			Summary: "List users", Tags: []string{"user"},
			Request: transport.Paginate{}, Response: transport.Envelope[transport.Page[*biz.User]]{}},
		{Method: fiber.MethodPost, Path: "/users", Handler: user.Create,
			Summary: "Create a user", Tags: []string{"user"},
			Request: UserAdd{}, Response: transport.Envelope[biz.User]{}},
		{Method: fiber.MethodGet, Path: "/users/:id", Handler: user.Get,
			Summary: "Get a user", Tags: []string{"user"},
			Request: UserID{}, Response: transport.Envelope[biz.User]{}},
		{Method: fiber.MethodPut, Path: "/users/:id", Handler: user.Update,
			Summary: "Update a user", Tags: []string{"user"},
			Request: UserUpdate{}, Response: transport.Envelope[biz.User]{}},
		{Method: fiber.MethodDelete, Path: "/users/:id", Handler: user.Delete,
			Summary: "Delete a user", Tags: []string{"user"},
			Request: UserID{}},
	}, nil
}
