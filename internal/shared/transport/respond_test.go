package transport_test

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/go-rio/rio"
	"github.com/gofiber/fiber/v3"
	"github.com/libtnb/assert/must"

	"github.com/libtnb/fiber-skeleton/internal/shared/apperr"
	"github.com/libtnb/fiber-skeleton/internal/shared/transport"
)

func respond(t *testing.T, err error) (int, string) {
	t.Helper()
	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error { return transport.ErrorFrom(c, err) })

	resp, aerr := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	must.NoError(t, aerr)
	body, aerr := io.ReadAll(resp.Body)
	must.NoError(t, aerr)
	return resp.StatusCode, string(body)
}

func TestErrorFromNotFound(t *testing.T) {
	status, body := respond(t, rio.ErrNotFound)
	must.Equal(t, status, fiber.StatusNotFound)
	must.Contains(t, body, "not found")
}

func TestErrorFromKinds(t *testing.T) {
	for kind, want := range map[apperr.Kind]int{
		apperr.KindInvalid:       fiber.StatusBadRequest,
		apperr.KindUnauthorized:  fiber.StatusUnauthorized,
		apperr.KindForbidden:     fiber.StatusForbidden,
		apperr.KindNotFound:      fiber.StatusNotFound,
		apperr.KindConflict:      fiber.StatusConflict,
		apperr.KindUnprocessable: fiber.StatusUnprocessableEntity,
	} {
		err := apperr.New(kind, "mod.code", "public detail").Errorf("internal detail")
		status, body := respond(t, err)
		must.Equal(t, status, want, must.Msgf("kind %s", kind))
		must.Contains(t, body, "mod.code")
		must.Contains(t, body, "public detail")
		must.NotContains(t, body, "internal detail")
	}
}

func TestErrorFromUnknownErrorHidesDetails(t *testing.T) {
	status, body := respond(t, errors.New("password=hunter2 exploded"))
	must.Equal(t, status, fiber.StatusInternalServerError)
	must.NotContains(t, body, "hunter2")
}
