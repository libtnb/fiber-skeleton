package bootstrap

import (
	"context"
	"log/slog"

	"github.com/go-rio/rio"
	"github.com/go-rio/sqlite"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/conf"
)

// Data owns the database handle for shutdown and /readyz; modules inject the
// plain *rio.DB instead.
type Data struct {
	DB *rio.DB
}

// NewData opens the database (swap SQLite for MySQL/PostgreSQL freely).
func NewData(i do.Injector) (*Data, error) {
	config := do.MustInvoke[*conf.Config](i)
	log := do.MustInvoke[*slog.Logger](i)

	db, err := sqlite.Open(
		"file:"+config.Database.Path+"?_txlock=immediate&_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)",
		rio.WithQueryHook(newSlogHook(log, config.Database.Debug)),
	)
	if err != nil {
		return nil, err
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

	return &Data{DB: db}, nil
}

// ProvideDB exposes the plain handle for the data layers.
func ProvideDB(i do.Injector) (*rio.DB, error) {
	return do.MustInvoke[*Data](i).DB, nil
}

func (d *Data) Shutdown() error {
	return d.DB.Close()
}

func (d *Data) HealthCheck(ctx context.Context) error {
	return d.DB.Unwrap().PingContext(ctx)
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
