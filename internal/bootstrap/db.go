package bootstrap

import (
	"log/slog"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/knadh/koanf/v2"
	"github.com/libtnb/sqlite"
	sloggorm "github.com/orandin/slog-gorm"
	"gorm.io/gorm"

	"github.com/libtnb/fiber-skeleton/internal/migration"
)

func NewDB(conf *koanf.Koanf, log *slog.Logger) (*gorm.DB, error) {
	// You can use any other database, like MySQL or PostgreSQL.
	return gorm.Open(sqlite.Open("file:"+conf.MustString("database.path")+"?_txlock=immediate&_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)"), &gorm.Config{
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
