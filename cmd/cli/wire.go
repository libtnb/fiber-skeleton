//go:build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/libtnb/fiber-skeleton/internal/app"
	"github.com/libtnb/fiber-skeleton/internal/bootstrap"
	"github.com/libtnb/fiber-skeleton/internal/data"
	"github.com/libtnb/fiber-skeleton/internal/route"
	"github.com/libtnb/fiber-skeleton/internal/service"
)

// initCli init command line.
func initCli() (*app.Cli, error) {
	panic(wire.Build(bootstrap.ProviderSet, route.ProviderSet, service.ProviderSet, data.ProviderSet, app.NewCli))
}
