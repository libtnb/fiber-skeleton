package app

import (
	"time"

	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/biz"
	"github.com/libtnb/fiber-skeleton/internal/bootstrap"
	"github.com/libtnb/fiber-skeleton/internal/command"
	"github.com/libtnb/fiber-skeleton/internal/data"
	"github.com/libtnb/fiber-skeleton/internal/job"
	"github.com/libtnb/fiber-skeleton/internal/route"
	"github.com/libtnb/fiber-skeleton/internal/service"
)

// NewInjector assembles every package of the application. Services are lazy:
// only what the invoked entry point actually needs gets built.
func NewInjector() do.Injector {
	return do.NewWithOpts(&do.InjectorOpts{
		// keeps /readyz bounded even when a dependency hangs
		HealthCheckGlobalTimeout: 5 * time.Second,
	},
		bootstrap.Package,
		biz.Package,
		data.Package,
		service.Package,
		route.Package,
		command.Package,
		job.Package,
		do.Lazy(NewApp),
		do.Lazy(NewCli),
	)
}
