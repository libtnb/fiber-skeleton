package main

import (
	"errors"
	"fmt"
	"os"

	_ "time/tzdata"

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

func run() (err error) {
	// keep stdout for command output: logs default to the file only
	if os.Getenv("APP_LOG__OUTPUT") == "" {
		_ = os.Setenv("APP_LOG__OUTPUT", "file")
	}

	cli, cleanup, err := app.InitializeCLI()
	if err != nil {
		return err
	}
	defer func() {
		if cleanup != nil {
			err = errors.Join(err, cleanup())
		}
	}()

	return cli.Run(version)
}
