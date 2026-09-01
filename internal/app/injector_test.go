package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/libtnb/assert/must"
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
	must.NoError(t, err)
	must.NotNil(t, application)

	// Every migration must compile to SQLite and run on an empty schema.
	must.NoError(t, application.migrator.Up(t.Context()))

	resp, err := application.router.Test(httptest.NewRequest(http.MethodGet, "/openapi.json", nil))
	must.NoError(t, err)
	must.Equal(t, resp.StatusCode, 200)
	body, err := io.ReadAll(resp.Body)
	must.NoError(t, err)
	must.NoError(t, resp.Body.Close())
	must.Contains(t, string(body), `"version": "test"`)
	must.Contains(t, string(body), `"/users/{id}"`)

	must.NoError(t, cleanup())
	must.NoError(t, cleanup(), "generated cleanup must be idempotent")

	management, cleanupCLI, err := InitializeCLI()
	must.NoError(t, err)
	must.NotNil(t, management)
	must.NoError(t, cleanupCLI())
	must.NoError(t, cleanupCLI(), "generated cleanup must be idempotent")
}
