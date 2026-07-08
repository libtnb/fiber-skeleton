package request

import (
	"github.com/gofiber/fiber/v3"
)

// WithPrepare runs after binding and before validation: fill defaults or
// normalize values. Authorization belongs in middleware or usecases.
type WithPrepare interface {
	Prepare(c fiber.Ctx) error
}

// WithRules ANDs extra rules onto the struct tags at runtime.
type WithRules interface {
	Rules(c fiber.Ctx) map[string]string
}

// WithFilters applies value filters (trim, lower, ...) before validation.
type WithFilters interface {
	Filters(c fiber.Ctx) map[string]string
}

// WithMessages overrides message templates for this request only.
type WithMessages interface {
	Messages(c fiber.Ctx) map[string]string
}
