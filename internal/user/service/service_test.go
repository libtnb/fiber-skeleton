package service_test

import (
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rio/rio"
	"github.com/gofiber/fiber/v3"
	"github.com/libtnb/assert/check"
	"github.com/libtnb/assert/must"
	"github.com/libtnb/validator"

	mocksbiz "github.com/libtnb/fiber-skeleton/internal/mocks/user/biz"
	"github.com/libtnb/fiber-skeleton/internal/user/biz"
	"github.com/libtnb/fiber-skeleton/internal/user/service"
)

// newTestApp wires the service against a mocked repo and a real validator.
func newTestApp(t *testing.T) (*fiber.App, *mocksbiz.UserRepoMock) {
	t.Helper()

	repo := &mocksbiz.UserRepoMock{}
	user := service.NewUserService(biz.NewUserUsecase(repo), validator.MustNew())

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
	repo.ListFunc = func(context.Context, int, int) ([]*biz.User, int64, error) {
		return []*biz.User{{ID: 1, Name: "alice"}}, 1, nil
	}

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/users", nil))

	must.NoError(t, err)
	check.Equal(t, resp.StatusCode, fiber.StatusOK)

	b, err := io.ReadAll(resp.Body)
	must.NoError(t, err)
	check.Contains(t, string(b), "alice")

	listed := repo.ListCalls()
	must.Len(t, listed, 1)
	// the paginate defaults reach the usecase
	check.Equal(t, listed[0].Page, 1)
	check.Equal(t, listed[0].Limit, 10)
}

func TestUserGet(t *testing.T) {
	app, repo := newTestApp(t)
	repo.GetFunc = func(context.Context, uint) (*biz.User, error) {
		return &biz.User{ID: 1, Name: "alice"}, nil
	}

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/users/1", nil))

	must.NoError(t, err)
	check.Equal(t, resp.StatusCode, fiber.StatusOK)
	got := repo.GetCalls()
	must.Len(t, got, 1)
	check.Equal(t, got[0].ID, uint(1))
}

func TestUserGet_NotFoundMapsTo404(t *testing.T) {
	app, repo := newTestApp(t)
	repo.GetFunc = func(context.Context, uint) (*biz.User, error) {
		return nil, rio.ErrNotFound
	}

	resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/users/9", nil))

	must.NoError(t, err)
	check.Equal(t, resp.StatusCode, fiber.StatusNotFound)
}

func TestUserCreate(t *testing.T) {
	app, repo := newTestApp(t)
	repo.ExistsNameFunc = func(context.Context, string) (bool, error) { return false, nil }
	repo.CreateFunc = func(context.Context, *biz.User) error { return nil }

	req := httptest.NewRequest(fiber.MethodPost, "/users", strings.NewReader(`{"name":"alice"}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := app.Test(req)

	must.NoError(t, err)
	check.Equal(t, resp.StatusCode, fiber.StatusOK)
	created := repo.CreateCalls()
	must.Len(t, created, 1)
	check.Equal(t, created[0].User.Name, "alice")
}

func TestUserCreate_NameTakenMapsToConflict(t *testing.T) {
	app, repo := newTestApp(t)
	repo.ExistsNameFunc = func(context.Context, string) (bool, error) { return true, nil }

	req := httptest.NewRequest(fiber.MethodPost, "/users", strings.NewReader(`{"name":"alice"}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := app.Test(req)

	must.NoError(t, err)
	check.Equal(t, resp.StatusCode, fiber.StatusConflict)
}

func TestUserCreate_RejectsShortName(t *testing.T) {
	app, _ := newTestApp(t) // no repo funcs: validation must fail first

	req := httptest.NewRequest(fiber.MethodPost, "/users", strings.NewReader(`{"name":"ab"}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := app.Test(req)

	must.NoError(t, err)
	check.Equal(t, resp.StatusCode, fiber.StatusUnprocessableEntity)
}

func TestUserUpdate_NotFoundMapsTo404(t *testing.T) {
	app, repo := newTestApp(t)
	repo.UpdateFunc = func(context.Context, *biz.User) (*biz.User, error) {
		return nil, rio.ErrNotFound
	}

	req := httptest.NewRequest(fiber.MethodPut, "/users/9", strings.NewReader(`{"name":"alice"}`))
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	resp, err := app.Test(req)

	must.NoError(t, err)
	check.Equal(t, resp.StatusCode, fiber.StatusNotFound)
	updated := repo.UpdateCalls()
	must.Len(t, updated, 1)
	check.Equal(t, updated[0].User.ID, uint(9))
	check.Equal(t, updated[0].User.Name, "alice")
}

func TestUserDelete(t *testing.T) {
	app, repo := newTestApp(t)
	repo.DeleteFunc = func(context.Context, uint) error { return nil }

	resp, err := app.Test(httptest.NewRequest(fiber.MethodDelete, "/users/1", nil))

	must.NoError(t, err)
	check.Equal(t, resp.StatusCode, fiber.StatusOK)
	deleted := repo.DeleteCalls()
	must.Len(t, deleted, 1)
	check.Equal(t, deleted[0].ID, uint(1))
}
