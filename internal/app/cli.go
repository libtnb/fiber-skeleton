package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v3"

	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
)

type Cli struct {
	cmd *cli.Command
}

func NewCli(cmd *cli.Command) *Cli {
	return &Cli{cmd: cmd}
}

// Run executes the command; SIGINT/SIGTERM cancel the context handed to it.
func (r *Cli) Run(version string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r.cmd.Version = version

	return r.cmd.Run(ctx, os.Args)
}

// newRootCommand assembles every command contribution into the root CLI.
func newRootCommand(commands registry.Commands) *cli.Command {
	return &cli.Command{
		Name:     "cli",
		Usage:    "management commands",
		Commands: commands,
	}
}
