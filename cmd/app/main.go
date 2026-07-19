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

// Errors go to stderr: the app logger's writer is already closed here.
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("[APP] version", version)

	injector := app.NewInjector(version)
	// closes every shutdownable service in reverse dependency order
	defer func() { _ = injector.Shutdown() }()

	application, err := do.Invoke[*app.App](injector)
	if err != nil {
		return err
	}

	return application.Run()
}
