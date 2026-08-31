//go:build wireinject

package app

import (
	"github.com/libtnb/wire"

	"github.com/libtnb/fiber-skeleton/internal/migrations"
	"github.com/libtnb/fiber-skeleton/internal/order"
	"github.com/libtnb/fiber-skeleton/internal/platform/bootstrap"
	"github.com/libtnb/fiber-skeleton/internal/platform/conf"
	"github.com/libtnb/fiber-skeleton/internal/platform/server"
	"github.com/libtnb/fiber-skeleton/internal/shared/registry"
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
	Provide(migrations.Collection).
	Provide(server.NewVersion).
	Provide(bootstrap.NewCron).
	Provide(server.NewRouter).
	Provide(newRootCommand).
	Provide(NewApp).
	Provide(NewCli)

var InitializeApp = ApplicationModule.Injector[func(string) (*App, func() error, error)]()

var InitializeCLI = ApplicationModule.Injector[func() (*Cli, func() error, error)]()
