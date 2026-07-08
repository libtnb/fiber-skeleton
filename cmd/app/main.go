package main

import (
	"fmt"
	"os"

	_ "time/tzdata"

	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/app"
)

// version is injected at build time: -ldflags "-X main.version=v1.2.3".
var version = "dev"

// Errors go to stderr directly: slog.SetDefault redirects the log package
// to the app logger, whose writer is already closed by the deferred Shutdown.
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("[APP] version", version)

	injector := app.NewInjector()
	// closes the database, the log writer and every other shutdownable
	// service in reverse dependency order
	defer func() { _ = injector.Shutdown() }()

	application, err := do.Invoke[*app.App](injector)
	if err != nil {
		return err
	}

	return application.Run()
}
