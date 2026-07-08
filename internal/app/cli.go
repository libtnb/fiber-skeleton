package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"
)

type Cli struct {
	cmd *cli.Command
}

func NewCli(i do.Injector) (*Cli, error) {
	return &Cli{
		cmd: do.MustInvoke[*cli.Command](i),
	}, nil
}

// Run executes the command; SIGINT/SIGTERM cancel the context handed to it.
func (r *Cli) Run(version string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r.cmd.Version = version

	return r.cmd.Run(ctx, os.Args)
}
