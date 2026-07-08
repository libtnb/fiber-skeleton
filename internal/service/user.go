package service

import (
	"github.com/gofiber/fiber/v3"

	"github.com/libtnb/fiber-skeleton/internal/biz"
	"github.com/libtnb/fiber-skeleton/internal/request"
)

// UserService adapts HTTP to the user usecase: bind, validate, delegate, respond.
type UserService struct {
	user *biz.UserUsecase
}

func NewUserService(user *biz.UserUsecase) *UserService {
	return &UserService{
		user: user,
	}
}

func (r *UserService) List(c fiber.Ctx) error {
	req, err := Bind[request.Paginate](c)
	if err != nil {
		return Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	users, total, err := r.user.List(c.Context(), req.Page, req.Limit)
	if err != nil {
		return ErrorFrom(c, err)
	}

	return Success(c, Page[*biz.User]{
		Total: total,
		Items: users,
	})
}

func (r *UserService) Get(c fiber.Ctx) error {
	req, err := Bind[request.UserID](c)
	if err != nil {
		return Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	user, err := r.user.Get(c.Context(), req.ID)
	if err != nil {
		return ErrorFrom(c, err)
	}

	return Success(c, user)
}

func (r *UserService) Create(c fiber.Ctx) error {
	req, err := Bind[request.UserAdd](c)
	if err != nil {
		return Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	user, err := r.user.Create(c.Context(), req.Name)
	if err != nil {
		return ErrorFrom(c, err)
	}

	return Success(c, user)
}

func (r *UserService) Update(c fiber.Ctx) error {
	req, err := Bind[request.UserUpdate](c)
	if err != nil {
		return Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	user, err := r.user.Update(c.Context(), req.ID, req.Name)
	if err != nil {
		return ErrorFrom(c, err)
	}

	return Success(c, user)
}

func (r *UserService) Delete(c fiber.Ctx) error {
	req, err := Bind[request.UserID](c)
	if err != nil {
		return Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	if err = r.user.Delete(c.Context(), req.ID); err != nil {
		return ErrorFrom(c, err)
	}

	return Success(c, nil)
}
