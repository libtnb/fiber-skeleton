//go:build wireinject

package app

import (
	"github.com/libtnb/wire"

	"github.com/libtnb/fiber-skeleton/internal/bootstrap"
	"github.com/libtnb/fiber-skeleton/internal/conf"
	"github.com/libtnb/fiber-skeleton/internal/order"
	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
	"github.com/libtnb/fiber-skeleton/internal/server"
	"github.com/libtnb/fiber-skeleton/internal/user"
)

var ApplicationModule = wire.New().
	Multibind[registry.Routes]().
	Multibind[registry.Commands]().
	Multibind[registry.Jobs]().
	Multibind[registry.Subscriptions]().
	Multibind[registry.HealthChecks]().
	Include(
		bootstrap.Module,
		server.Module,
		user.Module,
		order.Module,
	).
	Provide(conf.Load).
	Provide(server.NewVersion).
	Provide(bootstrap.NewCron).
	Provide(server.NewRouter).
	Provide(newRootCommand).
	Provide(NewApp).
	Provide(NewCli)

var InitializeApp = ApplicationModule.Injector[func(string) (*App, func() error, error)]()

var InitializeCLI = ApplicationModule.Injector[func() (*Cli, func() error, error)]()
