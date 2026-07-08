package command

import (
	"context"
	"fmt"
	"time"

	"github.com/libtnb/migrate"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"
)

// MigrateCommand contributes `cli migrate` with status and rollback.
func MigrateCommand(i do.Injector) (*cli.Command, error) {
	return &cli.Command{
		Name:  "migrate",
		Usage: "apply pending database migrations",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			m, err := do.Invoke[*migrate.Migrator](i)
			if err != nil {
				return err
			}
			if err = m.Up(ctx); err != nil {
				return err
			}
			fmt.Println("database migrated")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "status",
				Usage: "list migrations and whether they are applied",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					m, err := do.Invoke[*migrate.Migrator](i)
					if err != nil {
						return err
					}
					statuses, err := m.Status(ctx)
					if err != nil {
						return err
					}
					for _, s := range statuses {
						state := "pending"
						switch {
						case s.Drifted:
							state = "drifted"
						case !s.Registered:
							state = "applied, missing from the collection"
						case s.Applied:
							state = "applied " + s.AppliedAt.Local().Format(time.DateTime)
						}
						fmt.Printf("%s\t%s\n", s.Name, state)
					}
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "roll back the most recently applied migrations",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "step", Value: 1, Usage: "how many migrations to undo"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					m, err := do.Invoke[*migrate.Migrator](i)
					if err != nil {
						return err
					}
					if err = m.Rollback(ctx, cmd.Int("step")); err != nil {
						return err
					}
					fmt.Println("rollback complete")
					return nil
				},
			},
		},
	}, nil
}
