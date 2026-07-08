package request

import (
	"github.com/gofiber/fiber/v3"
)

// Paginate caps page*limit well below integer overflow.
type Paginate struct {
	Page  uint `json:"page" form:"page" query:"page" validate:"number && min:1 && max:1000000"`
	Limit uint `json:"limit" form:"limit" query:"limit" validate:"number && min:1 && max:1000"`
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
