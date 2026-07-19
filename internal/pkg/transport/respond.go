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

	"github.com/libtnb/fiber-skeleton/internal/pkg/apperr"
)

// Envelope is the one response shape, typed so routes can document bodies.
type Envelope[T any] struct {
	Msg  string `json:"msg"`
	Code string `json:"code,omitempty"`
	Data T      `json:"data,omitempty"`
}

// Page is the typed payload of list responses.
type Page[T any] struct {
	Total int64 `json:"total"`
	Items []T   `json:"items"`
}

func Success[T any](c fiber.Ctx, data T) error {
	return c.JSON(&Envelope[T]{
		Msg:  "success",
		Data: data,
	})
}

func Error(c fiber.Ctx, code int, format string, args ...any) error {
	return c.Status(code).JSON(&Envelope[any]{
		Msg: fmt.Sprintf(format, args...),
	})
}

// ErrorSystem writes a generic 500 without leaking details.
func ErrorSystem(c fiber.Ctx) error {
	return c.Status(http.StatusInternalServerError).JSON(&Envelope[any]{
		Msg: http.StatusText(http.StatusInternalServerError),
	})
}

// ErrorFrom maps known errors to their status; anything else logs and 500s.
func ErrorFrom(c fiber.Ctx, err error) error {
	if errors.Is(err, rio.ErrNotFound) {
		return Error(c, fiber.StatusNotFound, "not found")
	}

	if status := statusFromKind(apperr.KindOf(err)); status != 0 {
		return c.Status(status).JSON(&Envelope[any]{
			Msg:  oops.GetPublic(err, http.StatusText(status)),
			Code: apperr.CodeOf(err),
		})
	}

	slog.ErrorContext(c.Context(), "request failed",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.Any("error", err),
	)
	return ErrorSystem(c)
}

// statusFromKind maps kinds to statuses; 0 = no kind, do not expose.
func statusFromKind(kind apperr.Kind) int {
	switch kind {
	case apperr.KindInvalid:
		return fiber.StatusBadRequest
	case apperr.KindUnauthorized:
		return fiber.StatusUnauthorized
	case apperr.KindForbidden:
		return fiber.StatusForbidden
	case apperr.KindNotFound:
		return fiber.StatusNotFound
	case apperr.KindConflict:
		return fiber.StatusConflict
	case apperr.KindUnprocessable:
		return fiber.StatusUnprocessableEntity
	default:
		return 0
	}
}
