//go:build wireinject

// Package bootstrap provides the boot-time infrastructure; business modules
// never import it.
package bootstrap

import (
	"log/slog"

	"github.com/go-rio/migrate"
	"github.com/go-rio/rio"
	"github.com/libtnb/utils/crypt"
	"github.com/libtnb/validator"
	"github.com/libtnb/wire"

	"github.com/libtnb/fiber-skeleton/internal/pkg/event"
	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
)

var Module = wire.New().
	Provide(NewLogger).
	Provide(NewData).
	Provide(ProvideDB).
	Provide(NewCrypter).
	Provide(NewValidator).
	Provide(NewBus).
	Provide(NewMigrate).
	Multibind[registry.HealthChecks]().
	Contribute[registry.HealthChecks](DatabaseHealthCheck).
	Multibind[registry.Jobs]().
	Contribute[registry.Jobs](Heartbeat).
	Multibind[registry.Commands]().
	Contribute[registry.Commands](MigrateCommand).
	Export[*slog.Logger]().
	Export[*rio.DB]().
	Export[crypt.Crypter]().
	Export[*validator.Validator]().
	Export[event.Bus]().
	Export[*migrate.Migrator]().
	Export[registry.HealthChecks]().
	Export[registry.Jobs]().
	Export[registry.Commands]()
