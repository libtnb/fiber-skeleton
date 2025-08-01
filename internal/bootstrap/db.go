package bootstrap

import (
	"log/slog"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/knadh/koanf/v2"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/gormlite"
	sloggorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm"

	"github.com/libtnb/fiber-skeleton/internal/migration"
)

func NewDB(conf *koanf.Koanf, log *slog.Logger) (*gorm.DB, error) {
	// You can use any other database, like MySQL or PostgreSQL.
	return gorm.Open(gormlite.Open(conf.MustString("database.path")), &gorm.Config{
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
