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
	// keep stdout for command output: logs default to the file only
	if os.Getenv("APP_LOG__OUTPUT") == "" {
		_ = os.Setenv("APP_LOG__OUTPUT", "file")
	}

	injector := app.NewInjector()
	defer func() { _ = injector.Shutdown() }()

	cli, err := do.Invoke[*app.Cli](injector)
	if err != nil {
		return err
	}

	return cli.Run(version)
}
