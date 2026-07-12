// Package transport holds the HTTP helpers shared by every module's service
// layer: request binding/validation, response envelopes and error mapping.
package transport

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-rio/rio"
	"github.com/gofiber/fiber/v3"
	"github.com/samber/oops"
)

type SuccessResponse struct {
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Envelope mirrors SuccessResponse with a typed payload; route declarations
// use it to document response bodies.
type Envelope[T any] struct {
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// Page is the typed payload of list responses.
type Page[T any] struct {
	Total int64 `json:"total"`
	Items []T   `json:"items"`
}

type ErrorResponse struct {
	Msg  string `json:"msg"`
	Code string `json:"code,omitempty"`
}

func Success(c fiber.Ctx, data any) error {
	return c.JSON(&SuccessResponse{
		Msg:  "success",
		Data: data,
	})
}

func Error(c fiber.Ctx, code int, format string, args ...any) error {
	return c.Status(code).JSON(&ErrorResponse{
		Msg: fmt.Sprintf(format, args...),
	})
}

// ErrorSystem writes a generic 500 without leaking details.
func ErrorSystem(c fiber.Ctx) error {
	return c.Status(http.StatusInternalServerError).JSON(&ErrorResponse{
		Msg: http.StatusText(http.StatusInternalServerError),
	})
}

// ErrorFrom maps an error to an HTTP response. A not-found becomes 404; an oops
// error whose Code is known becomes that status and returns the error's public
// message and code; anything else is logged with its full structured context
// (stack trace, domain, attributes) and answered as a 500 without leaking it.
func ErrorFrom(c fiber.Ctx, err error) error {
	if errors.Is(err, rio.ErrNotFound) {
		return Error(c, fiber.StatusNotFound, "not found")
	}

	if oopsErr, ok := oops.AsError[oops.OopsError](err); ok {
		code, _ := oopsErr.Code().(string)
		if status := statusFromCode(code); status != 0 {
			return c.Status(status).JSON(&ErrorResponse{
				Msg:  oops.GetPublic(err, http.StatusText(status)),
				Code: code,
			})
		}
	}

	slog.ErrorContext(c.Context(), "request failed",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.Any("error", err),
	)
	return ErrorSystem(c)
}

// statusFromCode maps an application error code to an HTTP status, or 0 when
// the code is unknown — an unknown code is an unexpected error, not a
// client-facing one. Add a case here when a module introduces a new code.
func statusFromCode(code string) int {
	switch code {
	case "user.name_taken":
		return fiber.StatusConflict
	case "order.user_not_found":
		return fiber.StatusUnprocessableEntity
	default:
		return 0
	}
}
