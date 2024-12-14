package request

import (
	"github.com/gofiber/fiber/v3"
)

type Paginate struct {
	Page  uint `json:"page" form:"page" query:"page" validate:"required|number|min:1"`
	Limit uint `json:"limit" form:"limit" query:"limit" validate:"required|number|min:1|max:1000"`
}

func (r *Paginate) Prepare(c fiber.Ctx) error {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Limit == 0 {
		r.Limit = 10
	}
	return nil
}
