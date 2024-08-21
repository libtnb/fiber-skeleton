package route

import (
	"github.com/gofiber/fiber/v3"

	"github.com/TheTNB/go-web-skeleton/internal/service"
)

func Http(r fiber.Router) {
	r.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	user := service.NewUserService()
	r.Get("users", user.List)
	r.Post("users", user.Create)
	r.Get("users/:id", user.Get)
	r.Put("users/:id", user.Update)
	r.Delete("users/:id", user.Delete)
}
