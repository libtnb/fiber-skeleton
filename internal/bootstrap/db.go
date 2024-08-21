package bootstrap

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/TheTNB/go-web-skeleton/internal/app"
	"github.com/TheTNB/go-web-skeleton/internal/migration"
)

func initOrm() {
	logLevel := logger.Error
	if app.Conf.Bool("app.debug") {
		logLevel = logger.Info
	}
	// You can use any other database, like MySQL or PostgreSQL.
	db, err := gorm.Open(sqlite.Open(app.Conf.MustString("database.path")), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}
	app.Orm = db
}

func runMigrate() {
	migrator := gormigrate.New(app.Orm, &gormigrate.Options{
		ValidateUnknownMigrations: true,
	}, migration.Migrations)
	if err := migrator.Migrate(); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}
}
