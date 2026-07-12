package conf

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/parsers/yaml"
	env "github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// APP_-prefixed env vars override file values; "__" separates nesting
// levels: APP_HTTP__READ_TIMEOUT=30s -> http.read_timeout.
const envPrefix = "APP_"

// exampleKey ships in config.example.yml and must never be used in production.
const exampleKey = "a-long-string-with-32-characters"

type Config struct {
	App      App      `koanf:"app"`
	HTTP     HTTP     `koanf:"http"`
	Log      Log      `koanf:"log"`
	Database Database `koanf:"database"`
}

type App struct {
	Name string `koanf:"name"`
	// Key is the 32-byte secret behind cookie encryption and the crypter.
	Key   string `koanf:"key"`
	Debug bool   `koanf:"debug"`
	// Locale selects validator messages (zh_Hans, zh_Hant, ja, ko, es, ru); default English.
	Locale string `koanf:"locale"`
}

type HTTP struct {
	Debug   bool   `koanf:"debug"`
	Address string `koanf:"address"`
	// DebugAddress serves pprof/expvar on a private port when set; never expose it.
	DebugAddress string `koanf:"debug_address"`
	// CorsOrigins allows cross-origin requests; empty = same-origin only.
	CorsOrigins []string `koanf:"cors_origins"`
	// Docs serves the OpenAPI document and UI at /openapi.json and /docs.
	Docs bool `koanf:"docs"`
	// BodyLimit is the maximum request body size in KB; zero means 4096.
	BodyLimit int `koanf:"body_limit"`
	// HeaderLimit is the maximum request header size in bytes.
	HeaderLimit       int           `koanf:"header_limit"`
	ReadTimeout       time.Duration `koanf:"read_timeout"`
	WriteTimeout      time.Duration `koanf:"write_timeout"`
	IdleTimeout       time.Duration `koanf:"idle_timeout"`
	ReduceMemoryUsage bool          `koanf:"reduce_memory_usage"`
}

type Log struct {
	// Level is one of debug | info | warn | error.
	Level string `koanf:"level"`
	// Output is one of file | stdout | both.
	Output string `koanf:"output"`
	Path   string `koanf:"path"`
}

type Database struct {
	Debug bool   `koanf:"debug"`
	Path  string `koanf:"path"`
	// Pool knobs; zero keeps the driver default.
	MaxOpenConns    int           `koanf:"max_open_conns"`
	MaxIdleConns    int           `koanf:"max_idle_conns"`
	ConnMaxLifetime time.Duration `koanf:"conn_max_lifetime"`
}

// Load reads $APP_CONFIG (default config/config.yml), applies APP_*
// overrides and validates the result.
func Load() (*Config, error) {
	path := os.Getenv("APP_CONFIG")
	if path == "" {
		path = "config/config.yml"
	}

	k := koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("load config %s: %w", path, err)
	}
	if err := k.Load(env.Provider(".", env.Opt{
		Prefix: envPrefix,
		TransformFunc: func(key, value string) (string, any) {
			if key == "APP_CONFIG" {
				return "", nil // selects the file above, not a config value
			}
			key = strings.ToLower(strings.TrimPrefix(key, envPrefix))
			return strings.ReplaceAll(key, "__", "."), value
		},
	}), nil); err != nil {
		return nil, fmt.Errorf("load env overrides: %w", err)
	}

	conf := &Config{}
	if err := k.UnmarshalWithConf("", conf, koanf.UnmarshalConf{
		Tag: "koanf",
		DecoderConfig: &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToSliceHookFunc(","), // env values for list fields split on comma
			),
			Result:           conf,
			WeaklyTypedInput: true,
		},
	}); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	conf.fillDefaults()
	if err := conf.check(); err != nil {
		return nil, err
	}

	return conf, nil
}

func (c *Config) fillDefaults() {
	if c.HTTP.BodyLimit <= 0 {
		c.HTTP.BodyLimit = 4096
	}
	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
	if c.Log.Output == "" {
		c.Log.Output = "file"
	}
	if c.Log.Path == "" {
		c.Log.Path = "storage/logs/app.log"
	}
}

func (c *Config) check() error {
	if len(c.App.Key) != 32 {
		return fmt.Errorf("app.key must be exactly 32 characters, got %d", len(c.App.Key))
	}
	if c.App.Key == exampleKey {
		fmt.Fprintln(os.Stderr, "[WARN] app.key is still the example value, generate your own before deploying")
	}
	if c.HTTP.Address == "" {
		return errors.New("http.address must not be empty")
	}
	switch c.Log.Level {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("log.level must be debug, info, warn or error, got %q", c.Log.Level)
	}
	switch c.Log.Output {
	case "file", "stdout", "both":
	default:
		return fmt.Errorf("log.output must be file, stdout or both, got %q", c.Log.Output)
	}

	return nil
}

// SlogLevel returns the parsed slog level; check guarantees it is valid.
func (l Log) SlogLevel() slog.Level {
	var level slog.Level
	_ = level.UnmarshalText([]byte(l.Level))
	return level
}
