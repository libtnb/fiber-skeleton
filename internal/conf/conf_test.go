package conf_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/libtnb/fiber-skeleton/internal/conf"
)

func writeConfig(t *testing.T, yaml string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yml")
	require.NoError(t, os.WriteFile(path, []byte(yaml), 0o600))
	t.Setenv("APP_CONFIG", path)
}

const minimal = `
app:
  name: "test-app"
  key: "a-long-string-with-32-characters"
http:
  address: ":3000"
`

func TestLoadAppliesDefaults(t *testing.T) {
	writeConfig(t, minimal)

	c, err := conf.Load()
	require.NoError(t, err)

	require.Equal(t, "test-app", c.App.Name)
	require.Equal(t, 4096, c.HTTP.BodyLimit)
	require.Equal(t, "info", c.Log.Level)
	require.Equal(t, "file", c.Log.Output)
	require.Equal(t, "storage/logs/app.log", c.Log.Path)
}

func TestLoadEnvOverrides(t *testing.T) {
	writeConfig(t, minimal)
	t.Setenv("APP_HTTP__ADDRESS", ":8080")
	t.Setenv("APP_HTTP__READ_TIMEOUT", "30s")
	t.Setenv("APP_HTTP__CORS_ORIGINS", "https://a.example,https://b.example")
	t.Setenv("APP_LOG__OUTPUT", "stdout")

	c, err := conf.Load()
	require.NoError(t, err)

	require.Equal(t, ":8080", c.HTTP.Address)
	require.Equal(t, 30*time.Second, c.HTTP.ReadTimeout)
	require.Equal(t, []string{"https://a.example", "https://b.example"}, c.HTTP.CorsOrigins)
	require.Equal(t, "stdout", c.Log.Output)
}

func TestLoadRejectsBadValues(t *testing.T) {
	for name, tc := range map[string]struct {
		yaml    string
		env     map[string]string
		wantErr string
	}{
		"short key": {
			yaml:    "app:\n  key: \"short\"\nhttp:\n  address: \":3000\"\n",
			wantErr: "app.key",
		},
		"missing address": {
			yaml:    "app:\n  key: \"a-long-string-with-32-characters\"\n",
			wantErr: "http.address",
		},
		"bad log level": {
			yaml:    minimal,
			env:     map[string]string{"APP_LOG__LEVEL": "verbose"},
			wantErr: "log.level",
		},
		"bad log output": {
			yaml:    minimal,
			env:     map[string]string{"APP_LOG__OUTPUT": "syslog"},
			wantErr: "log.output",
		},
	} {
		t.Run(name, func(t *testing.T) {
			writeConfig(t, tc.yaml)
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			_, err := conf.Load()
			require.ErrorContains(t, err, tc.wantErr)
		})
	}
}

func TestLoadMissingFileFails(t *testing.T) {
	t.Setenv("APP_CONFIG", filepath.Join(t.TempDir(), "absent.yml"))
	_, err := conf.Load()
	require.Error(t, err)
}
