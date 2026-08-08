package service

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
	"github.com/libtnb/fiber-skeleton/internal/user/biz"
)

func UserRoutes(user *UserService) transport.Endpoints {
	return transport.Endpoints{
		{Method: fiber.MethodGet, Path: "/users", Handler: user.List,
			Summary: "List users", Tags: []string{"user"},
			Document: transport.Describe[transport.Paginate, transport.Envelope[transport.Page[*biz.User]]](http.StatusOK)},
		{Method: fiber.MethodPost, Path: "/users", Handler: user.Create,
			Summary: "Create a user", Tags: []string{"user"},
			Document: transport.Describe[UserAdd, transport.Envelope[biz.User]](http.StatusOK)},
		{Method: fiber.MethodGet, Path: "/users/:id", Handler: user.Get,
			Summary: "Get a user", Tags: []string{"user"},
			Document: transport.Describe[UserID, transport.Envelope[biz.User]](http.StatusOK)},
		{Method: fiber.MethodPut, Path: "/users/:id", Handler: user.Update,
			Summary: "Update a user", Tags: []string{"user"},
			Document: transport.Describe[UserUpdate, transport.Envelope[biz.User]](http.StatusOK)},
		{Method: fiber.MethodDelete, Path: "/users/:id", Handler: user.Delete,
			Summary: "Delete a user", Tags: []string{"user"},
			Document: transport.DescribeNoBody[UserID](http.StatusOK)},
	}
}
