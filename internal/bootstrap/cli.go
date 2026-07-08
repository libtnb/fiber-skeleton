package bootstrap

import (
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"

	"github.com/libtnb/fiber-skeleton/internal/command"
)

func NewCli(i do.Injector) (*cli.Command, error) {
	commands, err := command.Commands(i)
	if err != nil {
		return nil, err
	}

	return &cli.Command{
		Name:     "cli",
		Usage:    "management commands",
		Commands: commands,
	}, nil
}
