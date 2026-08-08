package service

import (
	"net/http"

	"github.com/gofiber/fiber/v3"

	"github.com/libtnb/fiber-skeleton/internal/order/biz"
	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
)

func OrderRoutes(order *OrderService) transport.Endpoints {
	return transport.Endpoints{
		{Method: fiber.MethodGet, Path: "/orders", Handler: order.List,
			Summary: "List orders", Tags: []string{"order"},
			Document: transport.Describe[transport.Paginate, transport.Envelope[transport.Page[*biz.Order]]](http.StatusOK)},
		{Method: fiber.MethodPost, Path: "/orders", Handler: order.Create,
			Summary: "Place an order", Tags: []string{"order"},
			Document: transport.Describe[OrderCreate, transport.Envelope[biz.Order]](http.StatusOK)},
		{Method: fiber.MethodGet, Path: "/orders/:id", Handler: order.Get,
			Summary: "Get an order", Tags: []string{"order"},
			Document: transport.Describe[OrderID, transport.Envelope[biz.Order]](http.StatusOK)},
		{Method: fiber.MethodDelete, Path: "/orders/:id", Handler: order.Delete,
			Summary: "Delete an order", Tags: []string{"order"},
			Document: transport.DescribeNoBody[OrderID](http.StatusOK)},
	}
}
