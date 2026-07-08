package service

import (
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/registry"
)

// Package wires the service layer; constructors stay container-free where
// they don't need the injector itself.
var Package = do.Package(
	do.Lazy(NewHealthService),
	registry.Lazy(NewUserService),
)
