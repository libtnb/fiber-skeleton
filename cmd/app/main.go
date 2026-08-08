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
	fmt.Println("[APP] version", version)

	application, cleanup, err := app.InitializeApp(version)
	if err != nil {
		return err
	}
	defer func() {
		if cleanup != nil {
			err = errors.Join(err, cleanup())
		}
	}()

	return application.Run()
}
