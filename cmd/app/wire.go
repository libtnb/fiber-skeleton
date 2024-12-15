//go:build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/go-rat/fiber-skeleton/internal/app"
	"github.com/go-rat/fiber-skeleton/internal/bootstrap"
	"github.com/go-rat/fiber-skeleton/internal/data"
	"github.com/go-rat/fiber-skeleton/internal/http/middleware"
	"github.com/go-rat/fiber-skeleton/internal/route"
	"github.com/go-rat/fiber-skeleton/internal/service"
)

// initApp init application.
func initApp() (*app.App, error) {
	panic(wire.Build(bootstrap.ProviderSet, middleware.ProviderSet, route.ProviderSet, service.ProviderSet, data.ProviderSet, app.NewApp))
}
