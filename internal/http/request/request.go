package request

import (
	"github.com/gofiber/fiber/v2"
)

type WithAuthorize interface {
	Authorize(c *fiber.Ctx) error
}

type WithPrepare interface {
	Prepare(c *fiber.Ctx) error
}

type WithRules interface {
	Rules(c *fiber.Ctx) map[string]string
}

type WithFilters interface {
	Filters(c *fiber.Ctx) map[string]string
}

type WithMessages interface {
	Messages(c *fiber.Ctx) map[string]string
}
