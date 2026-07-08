package service

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"

	"github.com/libtnb/fiber-skeleton/internal/biz"
)

// SuccessResponse is the envelope for successful responses.
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

// ErrorResponse is the envelope for error responses.
type ErrorResponse struct {
	Msg string `json:"msg"`
}

// Success writes data in the success envelope.
func Success(c fiber.Ctx, data any) error {
	return c.JSON(&SuccessResponse{
		Msg:  "success",
		Data: data,
	})
}

// Error writes a formatted message with the given status code.
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

// ErrorFrom maps business errors to HTTP responses: not-found becomes 404,
// anything else is logged and answered as a 500 without leaking details.
func ErrorFrom(c fiber.Ctx, err error) error {
	if errors.Is(err, biz.ErrNotFound) {
		return Error(c, fiber.StatusNotFound, "%v", err)
	}

	slog.ErrorContext(c.Context(), "request failed",
		slog.String("method", c.Method()),
		slog.String("path", c.Path()),
		slog.Any("err", err),
	)
	return ErrorSystem(c)
}
