// Package command holds the CLI commands, contributed under the "commands:"
// naming convention and assembled into the root command at startup.
// Commands resolve their dependencies inside Action, not at construction:
// help and unrelated commands then run without config or a database.
package command

import (
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"

	"github.com/libtnb/fiber-skeleton/internal/registry"
)

// Prefix marks command contributions; assemblers and tests collect by it.
const Prefix = "commands:"

// Package lists every command contribution; add yours here.
var Package = do.Package(
	do.LazyNamed(Prefix+"migrate", MigrateCommand),
	do.LazyNamed(Prefix+"user", UserCommand),
)

// Commands collects every "commands:*" contribution.
func Commands(i do.Injector) ([]*cli.Command, error) {
	return registry.Collect[*cli.Command](i, Prefix)
}
