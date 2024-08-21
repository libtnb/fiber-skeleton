package request

import "github.com/gofiber/fiber/v3"

type Request[T any] interface {
	*T
	PrepareForValidation(c fiber.Ctx) error
}
