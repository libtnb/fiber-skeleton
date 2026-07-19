package transport_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/libtnb/validator"
	"github.com/stretchr/testify/require"

	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
)

type createReq struct {
	Name string `json:"name" validate:"required && min:3 && max:10"`
}

func bindOn[T any](t *testing.T, method, target, body, contentType string) (*T, int) {
	t.Helper()

	var bound *T
	app := fiber.New()
	app.All("/bind/:id?", func(c fiber.Ctx) error {
		req, err := transport.Bind[T](c, validator.NewValidator())
		if err != nil {
			return transport.Error(c, fiber.StatusUnprocessableEntity, "%v", err)
		}
		bound = req
		return transport.Success[any](c, nil)
	})

	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, target, reader)
	if contentType != "" {
		req.Header.Set(fiber.HeaderContentType, contentType)
	}
	resp, err := app.Test(req)
	require.NoError(t, err)
	return bound, resp.StatusCode
}

func TestBindBodyAndValidate(t *testing.T) {
	got, status := bindOn[createReq](t, fiber.MethodPost, "/bind", `{"name":"alice"}`, fiber.MIMEApplicationJSON)
	require.Equal(t, fiber.StatusOK, status)
	require.Equal(t, "alice", got.Name)
}

func TestBindRejectsInvalid(t *testing.T) {
	_, status := bindOn[createReq](t, fiber.MethodPost, "/bind", `{"name":"ab"}`, fiber.MIMEApplicationJSON)
	require.Equal(t, fiber.StatusUnprocessableEntity, status)
}

func TestBindRunsPrepareHook(t *testing.T) {
	got, status := bindOn[transport.Paginate](t, fiber.MethodGet, "/bind", "", "")
	require.Equal(t, fiber.StatusOK, status)
	require.Equal(t, 1, got.Page)
	require.Equal(t, 10, got.Limit)

	got, status = bindOn[transport.Paginate](t, fiber.MethodGet, "/bind?page=3&limit=50", "", "")
	require.Equal(t, fiber.StatusOK, status)
	require.Equal(t, 3, got.Page)
	require.Equal(t, 50, got.Limit)
}

func TestBindQueryOverLimitFailsValidation(t *testing.T) {
	_, status := bindOn[transport.Paginate](t, fiber.MethodGet, "/bind?limit=5000", "", "")
	require.Equal(t, fiber.StatusUnprocessableEntity, status)
}

type uriReq struct {
	ID uint `uri:"id" validate:"required && number"`
}

func TestBindURI(t *testing.T) {
	got, status := bindOn[uriReq](t, fiber.MethodGet, "/bind/42", "", "")
	require.Equal(t, fiber.StatusOK, status)
	require.EqualValues(t, 42, got.ID)
}
