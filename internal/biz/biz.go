package biz

import (
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/registry"
)

// Package wires the biz layer; constructors stay container-free.
var Package = do.Package(
	registry.Lazy(NewUserUsecase),
)
