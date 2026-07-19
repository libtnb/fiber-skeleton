package service_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rio/rio"
	"github.com/gofiber/fiber/v3"
	"github.com/libtnb/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/libtnb/fiber-skeleton/internal/user/biz"
	"github.com/libtnb/fiber-skeleton/internal/user/service"
	mocksbiz "github.com/libtnb/fiber-skeleton/mocks/user/biz"
)

// newTestApp wires the service against a mocked repo and a real validator.
func newTestApp(t *testing.T) (*fiber.App, *mocksbiz.UserRepo) {
	t.Helper()

	repo := mocksbiz.NewUserRepo(t)
	user := service.NewUserService(biz.NewUserUsecase(repo), validator.NewValidator())

	app := fiber.New()
	app.Get("/users", user.List)
	app.Post("/users", user.Create)
	app.Get("/users/:id", user.Get)
	app.Put("/users/:id", user.Update)
	app.Delete("/users/:id", user.Delete)

	return app, repo
}

func TestUserList(t *testing.T) {
	app, repo := newTestApp(t)
	repo.EXPECT().List(mock.Anything, 1, 10).
		Return([]*biz.User{{ID: 1, Name: "alice"}}, int64(1), nil)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/users", nil))

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.Contains(t, string(b), "alice")
}

func TestUserGet(t *testing.T) {
	app, repo := newTestApp(t)
	repo.EXPECT().Get(mock.Anything, uint(1)).
		Return(&biz.User{ID: 1, Name: "alice"}, nil)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/users/1", nil))

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestUserGet_NotFoundMapsTo404(t *testing.T) {
	app, repo := newTestApp(t)
	repo.EXPECT().Get(mock.Anything, uint(9)).
		Return(nil, rio.ErrNotFound)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/users/9", nil))

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestUserCreate(t *testing.T) {
	app, repo := newTestApp(t)
	repo.EXPECT().ExistsName(mock.Anything, "alice").Return(false, nil)
	repo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(u *biz.User) bool {
		return u.Name == "alice"
	})).Return(nil)

	req := httptest.NewRequest(fiber.MethodPost, "/users", strings.NewReader(`{"name":"alice"}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestUserCreate_NameTakenMapsToConflict(t *testing.T) {
	app, repo := newTestApp(t)
	repo.EXPECT().ExistsName(mock.Anything, "alice").Return(true, nil)

	req := httptest.NewRequest(fiber.MethodPost, "/users", strings.NewReader(`{"name":"alice"}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
}

func TestUserCreate_RejectsShortName(t *testing.T) {
	app, _ := newTestApp(t) // no repo expectations: validation must fail first

	req := httptest.NewRequest(fiber.MethodPost, "/users", strings.NewReader(`{"name":"ab"}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
}

func TestUserUpdate_NotFoundMapsTo404(t *testing.T) {
	app, repo := newTestApp(t)
	repo.EXPECT().Update(mock.Anything, mock.MatchedBy(func(u *biz.User) bool {
		return u.ID == 9 && u.Name == "alice"
	})).Return(nil, rio.ErrNotFound)

	req := httptest.NewRequest(fiber.MethodPut, "/users/9", strings.NewReader(`{"name":"alice"}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
}

func TestUserDelete(t *testing.T) {
	app, repo := newTestApp(t)
	repo.EXPECT().Delete(mock.Anything, uint(1)).Return(nil)

	resp, err := app.Test(httptest.NewRequest(fiber.MethodDelete, "/users/1", nil))

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
