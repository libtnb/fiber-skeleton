package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/cloudflare/tableflip"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"github.com/knadh/koanf/v2"
	"github.com/robfig/cron/v3"
)

type App struct {
	conf     *koanf.Koanf
	router   *fiber.App
	migrator *gormigrate.Gormigrate
	cron     *cron.Cron
}

func NewApp(conf *koanf.Koanf, router *fiber.App, migrator *gormigrate.Gormigrate, cron *cron.Cron, _ *validate.Validation) *App {
	return &App{
		conf:     conf,
		router:   router,
		migrator: migrator,
		cron:     cron,
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
	if runtime.GOOS != "windows" {
		return r.runServer()
	}

	return r.runServerFallback()
}

// runServer graceful run server
func (r *App) runServer() error {
	upg, err := tableflip.New(tableflip.Options{})
	if err != nil {
		return err
	}
	defer upg.Stop()

	// By prefixing PID to log, easy to interrupt from another process.
	log.SetPrefix(fmt.Sprintf("[PID %d]", os.Getpid()))

	// Listen for the process signal to trigger the tableflip upgrade.
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP)
		for range sig {
			if err = upg.Upgrade(); err != nil {
				log.Println("[Graceful] upgrade failed:", err)
			}
		}
	}()

	fmt.Println("[HTTP] listening and serving on", r.conf.MustString("http.address"))
	ln, err := upg.Listen("tcp", r.conf.MustString("http.address"))
	if err != nil {
		return err
	}
	defer ln.Close()

	go func() {
		if err = r.router.Listener(ln); err != nil {
			log.Println("[HTTP] server error:", err)
		}
	}()

	// tableflip ready
	if err = upg.Ready(); err != nil {
		return err
	}

	fmt.Println("[Graceful] ready for upgrade")
	<-upg.Exit()

	// Make sure to set a deadline on exiting the process
	// after upg.Exit() is closed. No new upgrades can be
	// performed if the parent doesn't exit.
	return r.router.ShutdownWithTimeout(60 * time.Second)
}

// runServerFallback fallback for windows
func (r *App) runServerFallback() error {
	fmt.Println("[HTTP] listening and serving on", r.conf.MustString("http.address"))
	return r.router.Listen(r.conf.MustString("http.address"))
}
