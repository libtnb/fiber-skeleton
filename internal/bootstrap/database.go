package bootstrap

import (
	"context"
	"log/slog"

	"github.com/go-rio/rio"
	"github.com/go-rio/sqlite"

	"github.com/libtnb/fiber-skeleton/internal/conf"
	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
)

// Data owns the database handle for shutdown and /readyz; modules inject the
// plain *rio.DB instead.
type Data struct {
	DB *rio.DB
}

// NewData opens the database (swap SQLite for MySQL/PostgreSQL freely).
func NewData(config *conf.Config, log *slog.Logger) (*Data, func() error, error) {
	db, err := sqlite.Open(
		"file:"+config.Database.Path+"?_txlock=immediate&_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)",
		rio.WithQueryHook(newSlogHook(log, config.Database.Debug)),
	)
	if err != nil {
		return nil, nil, err
	}

	sqlDB := db.Unwrap()
	if config.Database.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.Database.MaxOpenConns)
	}
	if config.Database.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.Database.MaxIdleConns)
	}
	if config.Database.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(config.Database.ConnMaxLifetime)
	}

	return &Data{DB: db}, db.Close, nil
}

// ProvideDB exposes the plain handle for the data layers.
func ProvideDB(data *Data) *rio.DB {
	return data.DB
}

func (d *Data) HealthCheck(ctx context.Context) error {
	return d.DB.Unwrap().PingContext(ctx)
}

func DatabaseHealthCheck(data *Data) registry.HealthCheck {
	return registry.HealthCheck{Name: "database", Check: data.HealthCheck}
}

// slogHook logs failed statements always, all statements when debug is on.
type slogHook struct {
	log     *slog.Logger
	verbose bool
}

func newSlogHook(log *slog.Logger, verbose bool) rio.QueryHook {
	return slogHook{log: log, verbose: verbose}
}

func (h slogHook) BeforeQuery(ctx context.Context, _ *rio.QueryEvent) context.Context {
	return ctx
}

func (h slogHook) AfterQuery(ctx context.Context, e *rio.QueryEvent) {
	switch {
	case e.Err != nil:
		h.log.ErrorContext(ctx, "query failed",
			slog.String("op", e.Op),
			slog.String("query", e.Query),
			slog.Duration("elapsed", e.Duration),
			slog.Any("err", e.Err),
		)
	case h.verbose:
		h.log.DebugContext(ctx, "query",
			slog.String("op", e.Op),
			slog.String("query", e.Query),
			slog.Int64("rows", e.RowsAffected),
			slog.Duration("elapsed", e.Duration),
		)
	}
}
