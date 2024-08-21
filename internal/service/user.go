package service

import (
	"github.com/gofiber/fiber/v3"

	"github.com/TheTNB/go-web-skeleton/internal/biz"
	"github.com/TheTNB/go-web-skeleton/internal/data"
	"github.com/TheTNB/go-web-skeleton/internal/http/request"
)

type UserService struct {
	repo biz.UserRepo
}

func NewUserService() *UserService {
	return &UserService{
		repo: data.NewUserRepo(),
	}
}

func (r *UserService) List(c fiber.Ctx) error {
	req, err := Bind[request.Paginate](c)
	if err != nil {
		return Error(c, fiber.StatusUnprocessableEntity, err.Error())
	}
	users, total, err := r.repo.List(req.Page, req.Limit)
	if err != nil {
		return ErrorSystem(c)
	}

	return Success(c, map[string]any{
		"total": total,
		"items": users,
	})
}

func (r *UserService) Get(c fiber.Ctx) error {
	req, err := Bind[request.UserID](c)
	if err != nil {
		return Error(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	user, err := r.repo.Get(req.ID)
	if err != nil {
		return ErrorSystem(c)
	}

	return Success(c, user)
}

func (r *UserService) Create(c fiber.Ctx) error {
	req, err := Bind[request.AddUser](c)
	if err != nil {
		return Error(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	user := new(biz.User)
	user.Name = req.Name
	if err = r.repo.Save(user); err != nil {
		return ErrorSystem(c)
	}

	return Success(c, user)
}

func (r *UserService) Update(c fiber.Ctx) error {
	req, err := Bind[request.UpdateUser](c)
	if err != nil {
		return Error(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	user := new(biz.User)
	user.ID = req.ID
	user.Name = req.Name
	if err = r.repo.Save(user); err != nil {
		return ErrorSystem(c)
	}

	return Success(c, user)
}

func (r *UserService) Delete(c fiber.Ctx) error {
	req, err := Bind[request.UserID](c)
	if err != nil {
		return Error(c, fiber.StatusUnprocessableEntity, err.Error())
	}

	if err = r.repo.Delete(req.ID); err != nil {
		return ErrorSystem(c)
	}

	return Success(c, nil)
}
