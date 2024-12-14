package app

import (
	"log/slog"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/knadh/koanf/v2"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type App struct {
	conf     *koanf.Koanf
	router   *fiber.App
	db       *gorm.DB
	migrator *gormigrate.Gormigrate
	cron     *cron.Cron
	log      *slog.Logger
}

func NewApp(conf *koanf.Koanf, router *fiber.App, db *gorm.DB, migrator *gormigrate.Gormigrate, cron *cron.Cron, log *slog.Logger) *App {
	return &App{
		conf:     conf,
		router:   router,
		db:       db,
		migrator: migrator,
		cron:     cron,
		log:      log,
	}
}

func (r *App) Run() error {
	// migrate database
	if err := r.migrator.Migrate(); err != nil {
		return err
	}

	// start cron
	r.cron.Start()

	// start http server
	return r.router.Listen(r.conf.MustString("http.address"), fiber.ListenConfig{
		ListenerNetwork:       "tcp",
		EnablePrefork:         r.conf.Bool("http.prefork"),
		EnablePrintRoutes:     r.conf.Bool("http.debug"),
		DisableStartupMessage: !r.conf.Bool("http.debug"),
	})
}
