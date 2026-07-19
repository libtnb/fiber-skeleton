package app

import (
	"time"

	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/bootstrap"
	"github.com/libtnb/fiber-skeleton/internal/conf"
	"github.com/libtnb/fiber-skeleton/internal/order"
	"github.com/libtnb/fiber-skeleton/internal/server"
	"github.com/libtnb/fiber-skeleton/internal/user"
)

// NewInjector assembles every package of the application.
func NewInjector(version string) do.Injector {
	i := do.NewWithOpts(&do.InjectorOpts{
		// keeps /readyz bounded even when a dependency hangs
		HealthCheckGlobalTimeout: 5 * time.Second,
	},
		do.Lazy(func(i do.Injector) (*conf.Config, error) { return conf.Load() }),

		// boot-time infrastructure and the HTTP server
		bootstrap.Package,
		server.Package,

		// business modules
		user.Package,
		order.Package,

		// application lifecycle
		do.Lazy(newRootCommand),
		do.Lazy(NewApp),
		do.Lazy(NewCli),
	)

	do.ProvideValue(i, server.Version(version))
	return i
}
