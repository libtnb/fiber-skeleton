package command

import (
	"context"
	"fmt"

	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"

	"github.com/libtnb/fiber-skeleton/internal/biz"
)

// UserCommand contributes `cli user ...`, reusing the same usecase as HTTP.
func UserCommand(i do.Injector) (*cli.Command, error) {
	return &cli.Command{
		Name:  "user",
		Usage: "manage users",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "list users",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					user, err := do.Invoke[*biz.UserUsecase](i)
					if err != nil {
						return err
					}
					users, total, err := user.List(ctx, 1, 100)
					if err != nil {
						return err
					}
					for _, u := range users {
						fmt.Printf("%d\t%s\t%s\n", u.ID, u.Name, u.CreatedAt.String())
					}
					fmt.Printf("%d user(s) in total\n", total)
					return nil
				},
			},
			{
				Name:      "add",
				Usage:     "add a user",
				ArgsUsage: "<name>",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					name := cmd.Args().First()
					if name == "" {
						return fmt.Errorf("usage: user add %s", cmd.ArgsUsage)
					}
					user, err := do.Invoke[*biz.UserUsecase](i)
					if err != nil {
						return err
					}
					u, err := user.Create(ctx, name)
					if err != nil {
						return err
					}
					fmt.Printf("created user #%d %s\n", u.ID, u.Name)
					return nil
				},
			},
		},
	}, nil
}
