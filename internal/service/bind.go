package service

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/libtnb/validator"

	"github.com/libtnb/fiber-skeleton/internal/request"
)

// Bind binds and validates the request against the validator installed via
// validator.SetDefault.
func Bind[T any](c fiber.Ctx) (*T, error) {
	v := validator.Default()

	req := new(T)

	// bind body, query and uri parameters
	switch c.Method() {
	case fiber.MethodPost, fiber.MethodPut, fiber.MethodPatch, fiber.MethodDelete:
		if c.Request().Header.ContentLength() > 0 {
			if err := c.Bind().Body(req); err != nil {
				return nil, err
			}
		}
	}
	if err := c.Bind().Query(req); err != nil {
		return nil, err
	}
	if err := c.Bind().URI(req); err != nil {
		return nil, err
	}

	// request hooks
	if hook, ok := any(req).(request.WithPrepare); ok {
		if err := hook.Prepare(c); err != nil {
			return nil, err
		}
	}

	// prepare the validation
	vd := v.Struct(req)
	if hook, ok := any(req).(request.WithRules); ok {
		for field, expr := range hook.Rules(c) {
			if err := vd.AddRules(field, expr); err != nil {
				return nil, err
			}
		}
	}
	if hook, ok := any(req).(request.WithFilters); ok {
		for field, filters := range hook.Filters(c) {
			if err := vd.AddFilters(field, filters); err != nil {
				return nil, err
			}
		}
	}
	if hook, ok := any(req).(request.WithMessages); ok {
		if messages := hook.Messages(c); messages != nil {
			if err := vd.AddMessages(messages); err != nil {
				return nil, err
			}
		}
	}

	// validate
	vd.Validate(c.Context())
	if vd.Fails() {
		return nil, errors.New(vd.Errors().One())
	}

	// write filtered values (trim, lower, ...) back into the request struct
	if err := vd.SafeBind(req); err != nil {
		return nil, err
	}

	return req, nil
}
