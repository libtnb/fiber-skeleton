// Package user is the user module's assembly.
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
	registry.Lazy2(service.NewUserService),
	do.LazyNamed(registry.RoutePrefix+"user", service.UserRoutes),
	do.LazyNamed(registry.CommandPrefix+"user", service.UserCommand),
)
