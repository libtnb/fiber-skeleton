package conf_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/libtnb/assert/must"

	"github.com/libtnb/fiber-skeleton/internal/platform/conf"
)

func writeConfig(t *testing.T, yaml string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yml")
	must.NoError(t, os.WriteFile(path, []byte(yaml), 0o600))
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
	must.NoError(t, err)

	must.Equal(t, c.App.Name, "test-app")
	must.Equal(t, c.HTTP.BodyLimit, 4096)
	must.Equal(t, c.Log.Level, "info")
	must.Equal(t, c.Log.Output, "file")
	must.Equal(t, c.Log.Path, "storage/logs/app.log")
}

func TestLoadEnvOverrides(t *testing.T) {
	writeConfig(t, minimal)
	t.Setenv("APP_HTTP__ADDRESS", ":8080")
	t.Setenv("APP_HTTP__READ_TIMEOUT", "30s")
	t.Setenv("APP_HTTP__CORS_ORIGINS", "https://a.example,https://b.example")
	t.Setenv("APP_LOG__OUTPUT", "stdout")

	c, err := conf.Load()
	must.NoError(t, err)

	must.Equal(t, c.HTTP.Address, ":8080")
	must.Equal(t, c.HTTP.ReadTimeout, 30*time.Second)
	must.DeepEqual(t, c.HTTP.CorsOrigins, []string{"https://a.example", "https://b.example"})
	must.Equal(t, c.Log.Output, "stdout")
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
			must.ErrorContains(t, err, tc.wantErr)
		})
	}
}

func TestLoadMissingFileFails(t *testing.T) {
	t.Setenv("APP_CONFIG", filepath.Join(t.TempDir(), "absent.yml"))
	_, err := conf.Load()
	must.Error(t, err)
}
