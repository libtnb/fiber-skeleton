package transport_test

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/go-rio/rio"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"

	"github.com/libtnb/fiber-skeleton/internal/pkg/apperr"
	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
)

func respond(t *testing.T, err error) (int, string) {
	t.Helper()
	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error { return transport.ErrorFrom(c, err) })

	resp, aerr := app.Test(httptest.NewRequest(fiber.MethodGet, "/", nil))
	require.NoError(t, aerr)
	body, aerr := io.ReadAll(resp.Body)
	require.NoError(t, aerr)
	return resp.StatusCode, string(body)
}

func TestErrorFromNotFound(t *testing.T) {
	status, body := respond(t, rio.ErrNotFound)
	require.Equal(t, fiber.StatusNotFound, status)
	require.Contains(t, body, "not found")
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
		require.Equal(t, want, status, "kind %s", kind)
		require.Contains(t, body, "mod.code")
		require.Contains(t, body, "public detail")
		require.NotContains(t, body, "internal detail")
	}
}

func TestErrorFromUnknownErrorHidesDetails(t *testing.T) {
	status, body := respond(t, errors.New("password=hunter2 exploded"))
	require.Equal(t, fiber.StatusInternalServerError, status)
	require.NotContains(t, body, "hunter2")
}
