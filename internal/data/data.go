package data

import (
	"context"
	"errors"
	"log/slog"

	"github.com/libtnb/sqlite"
	sloggorm "github.com/orandin/slog-gorm"
	"github.com/samber/do/v2"
	"gorm.io/gorm"

	"github.com/libtnb/fiber-skeleton/internal/biz"
	"github.com/libtnb/fiber-skeleton/internal/config"
)

// Package wires the data layer.
var Package = do.Package(
	do.Lazy(NewData),
	do.Lazy(NewUserRepo),
)

// Data owns the database handle; the container closes it on shutdown and
// pings it for /readyz.
type Data struct {
	DB *gorm.DB
}

// NewData opens the database (swap SQLite for MySQL/PostgreSQL freely).
func NewData(i do.Injector) (*Data, error) {
	conf := do.MustInvoke[*config.Config](i)
	log := do.MustInvoke[*slog.Logger](i)

	gormLogger := []sloggorm.Option{sloggorm.WithHandler(log.Handler())}
	if conf.Database.Debug {
		gormLogger = append(gormLogger, sloggorm.WithTraceAll())
	}

	db, err := gorm.Open(sqlite.Open("file:"+conf.Database.Path+"?_txlock=immediate&_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)"), &gorm.Config{
		Logger:                                   sloggorm.New(gormLogger...),
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if conf.Database.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(conf.Database.MaxOpenConns)
	}
	if conf.Database.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(conf.Database.MaxIdleConns)
	}
	if conf.Database.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(conf.Database.ConnMaxLifetime)
	}

	return &Data{DB: db}, nil
}

func (d *Data) Shutdown() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Data) HealthCheck(ctx context.Context) error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// wrapErr translates driver errors into biz sentinel errors.
func wrapErr(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return biz.ErrNotFound
	}
	return err
}
