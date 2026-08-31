//go:build wireinject

// Package order is the order module's assembly.
package order

import (
	"github.com/libtnb/wire"

	"github.com/libtnb/fiber-skeleton/internal/order/biz"
	"github.com/libtnb/fiber-skeleton/internal/order/data"
	"github.com/libtnb/fiber-skeleton/internal/order/service"
	"github.com/libtnb/fiber-skeleton/internal/shared/registry"
)

var Module = wire.New().
	Provide(data.NewOrderRepo).
	Provide(data.NewUsers).
	Provide(biz.NewOrderUsecase).
	Provide(service.NewOrderService).
	Multibind[registry.Routes]().
	Contribute[registry.Routes](service.OrderRoutes).
	Multibind[registry.Subscriptions]().
	Contribute[registry.Subscriptions](service.NewOrderPlacedLogger).
	Export[registry.Routes]().
	Export[registry.Subscriptions]()
