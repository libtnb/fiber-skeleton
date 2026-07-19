package transport

import (
	"github.com/gofiber/fiber/v3"
)

// Paginate caps page*limit well below integer overflow.
type Paginate struct {
	Page  int `json:"page" form:"page" query:"page" validate:"number && min:1 && max:1000000"`
	Limit int `json:"limit" form:"limit" query:"limit" validate:"number && min:1 && max:1000"`
}

// Prepare fills defaults before validation runs.
func (r *Paginate) Prepare(c fiber.Ctx) error {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Limit == 0 {
		r.Limit = 10
	}
	return nil
}

// WithPrepare runs between binding and validation: fill defaults, normalize.
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
