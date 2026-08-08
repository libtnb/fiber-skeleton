//go:build wireinject

// Package user is the user module's assembly.
package user

import (
	"github.com/libtnb/wire"

	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
	"github.com/libtnb/fiber-skeleton/internal/user/biz"
	"github.com/libtnb/fiber-skeleton/internal/user/data"
	"github.com/libtnb/fiber-skeleton/internal/user/service"
)

var Module = wire.New().
	Provide(data.NewUserRepo).
	Provide(biz.NewUserUsecase).
	Provide(service.NewUserService).
	Multibind[registry.Routes]().
	Contribute[registry.Routes](service.UserRoutes).
	Multibind[registry.Commands]().
	Contribute[registry.Commands](service.UserCommand).
	Export[*biz.UserUsecase]().
	Export[registry.Routes]().
	Export[registry.Commands]()
