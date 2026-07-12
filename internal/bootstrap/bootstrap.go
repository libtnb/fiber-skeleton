package bootstrap

import (
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
)

var Package = do.Package(
	do.Lazy(NewLogger),
	do.Lazy(NewSlog),
	do.Lazy(NewData),
	do.Lazy(ProvideDB),
	do.Lazy(NewCrypter),
	do.Lazy(NewValidator),
	do.Lazy(NewBus),
	do.Lazy(NewCron),
	do.LazyNamed(registry.JobPrefix+"heartbeat", Heartbeat),
	do.Lazy(NewMigrate),
	do.LazyNamed(registry.CommandPrefix+"migrate", MigrateCommand),
)
