// Package service adapts HTTP and CLI to the usecase: bind, validate,
// delegate, respond.
package service

import (
	"github.com/gofiber/fiber/v3"
	"github.com/libtnb/validator"

	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
	"github.com/libtnb/fiber-skeleton/internal/user/biz"
)

type UserService struct {
	user     *biz.UserUsecase
	validate *validator.Validator
}

func NewUserService(user *biz.UserUsecase, validate *validator.Validator) *UserService {
	return &UserService{
		user:     user,
		validate: validate,
	}
}

func (r *UserService) List(c fiber.Ctx) error {
	req, err := transport.Bind[transport.Paginate](c, r.validate)
	if err != nil {
		return transport.Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	users, total, err := r.user.List(c.Context(), req.Page, req.Limit)
	if err != nil {
		return transport.ErrorFrom(c, err)
	}

	return transport.Success(c, transport.Page[*biz.User]{
		Total: total,
		Items: users,
	})
}

func (r *UserService) Get(c fiber.Ctx) error {
	req, err := transport.Bind[UserID](c, r.validate)
	if err != nil {
		return transport.Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	user, err := r.user.Get(c.Context(), req.ID)
	if err != nil {
		return transport.ErrorFrom(c, err)
	}

	return transport.Success(c, user)
}

func (r *UserService) Create(c fiber.Ctx) error {
	req, err := transport.Bind[UserAdd](c, r.validate)
	if err != nil {
		return transport.Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	user, err := r.user.Create(c.Context(), req.Name)
	if err != nil {
		return transport.ErrorFrom(c, err)
	}

	return transport.Success(c, user)
}

func (r *UserService) Update(c fiber.Ctx) error {
	req, err := transport.Bind[UserUpdate](c, r.validate)
	if err != nil {
		return transport.Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	user, err := r.user.Update(c.Context(), req.ID, req.Name)
	if err != nil {
		return transport.ErrorFrom(c, err)
	}

	return transport.Success(c, user)
}

func (r *UserService) Delete(c fiber.Ctx) error {
	req, err := transport.Bind[UserID](c, r.validate)
	if err != nil {
		return transport.Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	if err = r.user.Delete(c.Context(), req.ID); err != nil {
		return transport.ErrorFrom(c, err)
	}

	return transport.Success[any](c, nil)
}
