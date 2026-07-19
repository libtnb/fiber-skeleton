// Package service adapts HTTP to the usecase: bind, validate, delegate,
// respond.
package service

import (
	"github.com/gofiber/fiber/v3"
	"github.com/libtnb/validator"

	"github.com/libtnb/fiber-skeleton/internal/order/biz"
	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
)

// OrderService adapts HTTP to the order usecase: bind, validate, delegate, respond.
type OrderService struct {
	order    *biz.OrderUsecase
	validate *validator.Validator
}

func NewOrderService(order *biz.OrderUsecase, validate *validator.Validator) *OrderService {
	return &OrderService{
		order:    order,
		validate: validate,
	}
}

func (r *OrderService) List(c fiber.Ctx) error {
	req, err := transport.Bind[transport.Paginate](c, r.validate)
	if err != nil {
		return transport.Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	orders, total, err := r.order.List(c.Context(), req.Page, req.Limit)
	if err != nil {
		return transport.ErrorFrom(c, err)
	}

	return transport.Success(c, transport.Page[*biz.Order]{
		Total: total,
		Items: orders,
	})
}

func (r *OrderService) Get(c fiber.Ctx) error {
	req, err := transport.Bind[OrderID](c, r.validate)
	if err != nil {
		return transport.Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	order, err := r.order.Get(c.Context(), req.ID)
	if err != nil {
		return transport.ErrorFrom(c, err)
	}

	return transport.Success(c, order)
}

func (r *OrderService) Create(c fiber.Ctx) error {
	req, err := transport.Bind[OrderCreate](c, r.validate)
	if err != nil {
		return transport.Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	order, err := r.order.Place(c.Context(), req.UserID, req.Amount)
	if err != nil {
		return transport.ErrorFrom(c, err)
	}

	return transport.Success(c, order)
}

func (r *OrderService) Delete(c fiber.Ctx) error {
	req, err := transport.Bind[OrderID](c, r.validate)
	if err != nil {
		return transport.Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	if err = r.order.Delete(c.Context(), req.ID); err != nil {
		return transport.ErrorFrom(c, err)
	}

	return transport.Success[any](c, nil)
}
