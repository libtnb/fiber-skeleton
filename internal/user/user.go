// Package user is the user module's assembly: the do.Package that wires its
// biz, data and service layers and registers its route and command contributions.
package user

import (
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
	"github.com/libtnb/fiber-skeleton/internal/user/biz"
	"github.com/libtnb/fiber-skeleton/internal/user/data"
	"github.com/libtnb/fiber-skeleton/internal/user/service"
)

var Package = do.Package(
	do.Lazy(data.NewUserRepo),
	registry.Lazy(biz.NewUserUsecase),
	registry.Lazy(service.NewUserService),
	do.LazyNamed(registry.RoutePrefix+"user", service.UserRoutes),
	do.LazyNamed(registry.CommandPrefix+"user", service.UserCommand),
)
