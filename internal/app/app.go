package app

import (
	"context"
	_ "expvar" // registers /debug/vars on the default mux
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof" // registers /debug/pprof on the default mux
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/libtnb/cron"
	"github.com/libtnb/graceful"
	"github.com/libtnb/migrate"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/config"
)

type App struct {
	conf     *config.Config
	router   *fiber.App
	migrator *migrate.Migrator
	cron     *cron.Cron
}

func NewApp(i do.Injector) (*App, error) {
	return &App{
		conf:     do.MustInvoke[*config.Config](i),
		router:   do.MustInvoke[*fiber.App](i),
		migrator: do.MustInvoke[*migrate.Migrator](i),
		cron:     do.MustInvoke[*cron.Cron](i),
	}, nil
}

// Run migrates the database, then hands the lifecycle to graceful:
// SIGINT/SIGTERM drains everything, SIGHUP hot-upgrades the binary.
func (r *App) Run() error {
	if err := r.migrator.Up(context.Background()); err != nil {
		return err
	}
	fmt.Println("[DB] database migrated")

	g := graceful.New(
		graceful.WithUpgrade(),
		graceful.WithShutdownTimeout(30*time.Second),
	)
	// pprof/expvar live on http.DefaultServeMux, served on a private port
	if addr := r.conf.HTTP.DebugAddress; addr != "" {
		g.Listen("debug", addr, &http.Server{})
	}
	g.Add("cron", r.cron.Start, r.cron.Stop)
	g.Listen("http", r.conf.HTTP.Address, fiberServer{app: r.router, conf: r.conf})

	fmt.Println("[HTTP] listening and serving on", r.conf.HTTP.Address)
	return g.Run()
}

// fiberServer adapts *fiber.App to graceful.Server.
type fiberServer struct {
	app  *fiber.App
	conf *config.Config
}

func (s fiberServer) Serve(ln net.Listener) error {
	return s.app.Listener(ln, fiber.ListenConfig{
		EnablePrintRoutes:     s.conf.HTTP.Debug,
		DisableStartupMessage: !s.conf.HTTP.Debug,
	})
}

func (s fiberServer) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
