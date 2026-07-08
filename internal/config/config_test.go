package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/libtnb/fiber-skeleton/internal/config"
)

func writeConfig(t *testing.T, content string) {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	t.Setenv("APP_CONFIG", path)
}

const validConfig = `
app:
  name: "test-app"
  key: "01234567890123456789012345678901"
http:
  address: ":3000"
`

func TestLoad(t *testing.T) {
	writeConfig(t, validConfig)

	conf, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, "test-app", conf.App.Name)
	// defaults are filled in
	assert.Equal(t, "info", conf.Log.Level)
	assert.Equal(t, "file", conf.Log.Output)
}

func TestLoad_EnvOverridesFile(t *testing.T) {
	writeConfig(t, validConfig)
	t.Setenv("APP_HTTP__ADDRESS", ":9999")
	t.Setenv("APP_LOG__LEVEL", "error")

	conf, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, ":9999", conf.HTTP.Address)
	assert.Equal(t, "error", conf.Log.Level)
}

func TestLoad_RejectsBadKey(t *testing.T) {
	writeConfig(t, `
app:
  key: "too-short"
http:
  address: ":3000"
`)

	_, err := config.Load()

	assert.ErrorContains(t, err, "app.key")
}

func TestLoad_RejectsBadLogOutput(t *testing.T) {
	writeConfig(t, validConfig+`
log:
  output: "nowhere"
`)

	_, err := config.Load()

	assert.ErrorContains(t, err, "log.output")
}

// TestExampleConfigLoads pins config.example.yml to the Config struct: if a
// new field or a rename breaks the example, this fails before a user does
// `make init` and hits it.
func TestExampleConfigLoads(t *testing.T) {
	t.Setenv("APP_CONFIG", "../../config/config.example.yml")

	_, err := config.Load()

	require.NoError(t, err)
}

func TestLoad_EnvListOverride(t *testing.T) {
	writeConfig(t, validConfig)
	t.Setenv("APP_HTTP__CORS_ORIGINS", "https://a.example,https://b.example")

	conf, err := config.Load()

	require.NoError(t, err)
	assert.Equal(t, []string{"https://a.example", "https://b.example"}, conf.HTTP.CorsOrigins)
}
