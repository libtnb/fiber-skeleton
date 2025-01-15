package bootstrap

import (
	"log/slog"

	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/knadh/koanf/v2"
	sloggorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm"

	"github.com/go-rat/fiber-skeleton/internal/migration"
)

func NewDB(conf *koanf.Koanf, log *slog.Logger) (*gorm.DB, error) {
	// You can use any other database, like MySQL or PostgreSQL.
	return gorm.Open(sqlite.Open(conf.MustString("database.path")), &gorm.Config{
		Logger:                                   sloggorm.New(sloggorm.WithHandler(log.Handler())),
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}

func NewMigrate(db *gorm.DB) *gormigrate.Gormigrate {
	return gormigrate.New(db, &gormigrate.Options{
		UseTransaction: true, // Note: MySQL not support DDL transaction
	}, migration.Migrations)
}
