package bootstrap

import (
	"log/slog"

	"github.com/libtnb/migrate"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/data"
	_ "github.com/libtnb/fiber-skeleton/internal/migration" // registers the migrations
)

func NewMigrate(i do.Injector) (*migrate.Migrator, error) {
	db, err := do.MustInvoke[*data.Data](i).DB.DB()
	if err != nil {
		return nil, err
	}

	return migrate.New(db, migrate.SQLite,
		migrate.WithLogger(do.MustInvoke[*slog.Logger](i)),
	)
}
