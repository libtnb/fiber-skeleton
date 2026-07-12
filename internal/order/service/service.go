// Package service adapts HTTP to the order usecase and owns the module's
// request DTOs, route contribution and event subscribers.
package service

import (
	"github.com/gofiber/fiber/v3"

	"github.com/libtnb/fiber-skeleton/internal/order/biz"
	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
)

// OrderService adapts HTTP to the order usecase: bind, validate, delegate, respond.
type OrderService struct {
	order *biz.OrderUsecase
}

func NewOrderService(order *biz.OrderUsecase) *OrderService {
	return &OrderService{
		order: order,
	}
}

func (r *OrderService) List(c fiber.Ctx) error {
	req, err := transport.Bind[transport.Paginate](c)
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
	req, err := transport.Bind[OrderID](c)
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
	req, err := transport.Bind[OrderCreate](c)
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
	req, err := transport.Bind[OrderID](c)
	if err != nil {
		return transport.Error(c, fiber.StatusUnprocessableEntity, "%v", err)
	}

	if err = r.order.Delete(c.Context(), req.ID); err != nil {
		return transport.ErrorFrom(c, err)
	}

	return transport.Success(c, nil)
}
