package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGeneratedGraphs builds both generated object graphs, catching wiring
// and resource-lifecycle mistakes early.
func TestGeneratedGraphs(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("APP_CONFIG", "../../config/config.example.yml")
	t.Setenv("APP_DATABASE__PATH", filepath.Join(tmp, "test.db"))
	t.Setenv("APP_LOG__OUTPUT", "file")
	t.Setenv("APP_LOG__PATH", filepath.Join(tmp, "test.log"))
	t.Setenv("APP_HTTP__DOCS", "true")

	application, cleanup, err := InitializeApp("test")
	require.NoError(t, err)
	require.NotNil(t, application)

	// Every migration must compile to SQLite and run on an empty schema.
	require.NoError(t, application.migrator.Up(t.Context()))

	resp, err := application.router.Test(httptest.NewRequest(http.MethodGet, "/openapi.json", nil))
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	require.Contains(t, string(body), `"version": "test"`)
	require.Contains(t, string(body), `"/users/{id}"`)

	require.NoError(t, cleanup())
	require.NoError(t, cleanup(), "generated cleanup must be idempotent")

	management, cleanupCLI, err := InitializeCLI()
	require.NoError(t, err)
	require.NotNil(t, management)
	require.NoError(t, cleanupCLI())
	require.NoError(t, cleanupCLI(), "generated cleanup must be idempotent")
}
