// Package order is the order module's assembly.
package order

import (
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/order/biz"
	"github.com/libtnb/fiber-skeleton/internal/order/data"
	"github.com/libtnb/fiber-skeleton/internal/order/service"
	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
)

var Package = do.Package(
	do.Lazy(data.NewOrderRepo),
	do.Lazy(data.NewUsers), // implements biz.Users over the user module
	registry.Lazy3(biz.NewOrderUsecase),
	registry.Lazy2(service.NewOrderService),
	do.LazyNamed(registry.RoutePrefix+"order", service.OrderRoutes),
	do.LazyNamed(registry.SubscriberPrefix+"order-placed", service.NewOrderPlacedLogger),
)
