package app

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
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
	fmt.Println("[DB] database migrated")

	// start cron scheduler
	r.cron.Start()
	fmt.Println("[CRON] cron scheduler started")

	// run http server
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		return r.runServerGraceful()
	}

	return r.runServer()
}

// runServer fallback for unsupported graceful OS
func (r *App) runServer() error {
	fmt.Println("[HTTP] Listening and serving HTTP on", r.conf.MustString("http.address"))
	return r.router.Listen(r.conf.MustString("http.address"), fiber.ListenConfig{
		ListenerNetwork:       "tcp",
		EnablePrefork:         r.conf.Bool("http.prefork"),
		EnablePrintRoutes:     r.conf.Bool("http.debug"),
		DisableStartupMessage: !r.conf.Bool("http.debug"),
	})
}

// runServerGraceful graceful for linux and darwin
func (r *App) runServerGraceful() error {
	upg, err := tableflip.New(tableflip.Options{})
	if err != nil {
		return err
	}
	defer upg.Stop()

	// By prefixing PID to log, easy to interrupt from another process.
	log.SetPrefix(fmt.Sprintf("[PID: %d]", os.Getpid()))

	// Listen for the process signal to trigger the tableflip upgrade.
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP)
		for range sig {
			if err = upg.Upgrade(); err != nil {
				log.Printf("graceful upgrade failed: %v", err)
			}
		}
	}()

	ln, err := upg.Listen("tcp", r.conf.MustString("http.address"))
	if err != nil {
		return err
	}
	defer ln.Close()

	fmt.Println("[HTTP] Listening and serving HTTP graceful on", r.conf.MustString("http.address"))
	go r.router.Listener(ln, fiber.ListenConfig{
		ListenerNetwork:       "tcp",
		EnablePrefork:         r.conf.Bool("http.prefork"),
		EnablePrintRoutes:     r.conf.Bool("http.debug"),
		DisableStartupMessage: !r.conf.Bool("http.debug"),
	})

	// tableflip ready
	if err = upg.Ready(); err != nil {
		return err
	}

	<-upg.Exit()

	// Make sure to set a deadline on exiting the process
	// after upg.Exit() is closed. No new upgrades can be
	// performed if the parent doesn't exit.
	return r.router.ShutdownWithTimeout(60 * time.Second)
}
