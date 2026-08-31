//go:build wireinject

package server

import (
	"github.com/libtnb/wire"

	"github.com/libtnb/fiber-skeleton/internal/shared/registry"
)

var Module = wire.New().
	Multibind[registry.Routes]().
	Contribute[registry.Routes](HealthRoutes).
	Contribute[registry.Routes](WsRoutes).
	Export[registry.Routes]()
