package route

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/go-rat/fiber-skeleton/internal/service"
)

type Cli struct {
	user *service.UserService
}

func NewCli(user *service.UserService) *Cli {
	return &Cli{
		user: user,
	}
}

func (r *Cli) Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "user",
			Usage: "用户",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return nil
			},
		},
	}
}
