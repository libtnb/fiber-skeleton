package request

import (
	"github.com/gofiber/fiber/v3"
)

type Paginate struct {
	Page  uint `json:"page" form:"page" query:"page" validate:"required,number,gte=1" comment:"页码"`
	Limit uint `json:"limit" form:"limit" query:"limit" validate:"required,number,gte=1,lte=1000" comment:"每页数量"`
}

func (r *Paginate) PrepareForValidation(c fiber.Ctx) error {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Limit == 0 {
		r.Limit = 10
	}
	return nil
}
