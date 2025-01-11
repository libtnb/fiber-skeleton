package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/go-rat/fiber-skeleton/internal/service"
)

type Http struct {
	user *service.UserService
}

func NewHttp(user *service.UserService) *Http {
	return &Http{
		user: user,
	}
}

func (r *Http) Register(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	router.Get("/users", r.user.List)
	router.Post("/users", r.user.Create)
	router.Get("/users/:id", r.user.Get)
	router.Put("/users/:id", r.user.Update)
	router.Delete("/users/:id", r.user.Delete)
}
