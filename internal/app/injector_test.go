package app_test

import (
	"path/filepath"
	"testing"

	"github.com/go-rio/migrate"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"

	"github.com/libtnb/fiber-skeleton/internal/app"
	"github.com/libtnb/fiber-skeleton/internal/bootstrap"
	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
	"github.com/libtnb/fiber-skeleton/internal/server"
)

// TestContainer builds the full object graph, catching wiring mistakes
// (missing providers, bad contributions) at test time instead of at startup.
func TestContainer(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("APP_CONFIG", "../../config/config.example.yml")
	t.Setenv("APP_DATABASE__PATH", filepath.Join(tmp, "test.db"))
	t.Setenv("APP_LOG__OUTPUT", "file")
	t.Setenv("APP_LOG__PATH", filepath.Join(tmp, "test.log"))

	injector := app.NewInjector()
	defer func() { _ = injector.Shutdown() }()

	_, err := do.Invoke[*app.App](injector)
	require.NoError(t, err)

	_, err = do.Invoke[*app.Cli](injector)
	require.NoError(t, err)

	// every named contribution must carry a known prefix: a typo like
	// "route:user" would otherwise be dropped silently
	require.NoError(t, registry.Verify(injector, registry.RoutePrefix, registry.CommandPrefix, registry.JobPrefix, registry.SubscriberPrefix))

	routes, err := registry.Collect[server.Endpoints](injector, registry.RoutePrefix)
	require.NoError(t, err)
	require.NotEmpty(t, routes)

	commands, err := registry.Collect[*cli.Command](injector, registry.CommandPrefix)
	require.NoError(t, err)
	require.NotEmpty(t, commands)

	jobs, err := registry.Collect[bootstrap.JobFn](injector, registry.JobPrefix)
	require.NoError(t, err)
	require.NotEmpty(t, jobs)

	// apply the migrations against the tmp database: every declaration must
	// compile to SQLite and run on an empty schema
	m, err := do.Invoke[*migrate.Migrator](injector)
	require.NoError(t, err)
	require.NoError(t, m.Up(t.Context()))
}
