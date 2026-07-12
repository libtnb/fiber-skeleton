// Package bootstrap holds the boot-time wiring: the providers that are built
// once at startup and assembled into the container — the logger, crypter,
// validator, scheduler and migrator. Business modules never import it; it is
// the composition layer beneath them.
package bootstrap

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/libtnb/logrotate"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/conf"
)

// Logger owns the rotating writer so the container can close it on shutdown;
// inject *slog.Logger everywhere else.
type Logger struct {
	*slog.Logger
	close func() error
}

func (l *Logger) Shutdown() error {
	return l.close()
}

// NewLogger builds the logger writing to a rotated file, stdout, or both.
func NewLogger(i do.Injector) (*Logger, error) {
	config := do.MustInvoke[*conf.Config](i)

	var (
		writers []io.Writer
		closer  = func() error { return nil }
	)

	if config.Log.Output == "file" || config.Log.Output == "both" {
		w, err := logrotate.New(config.Log.Path,
			logrotate.WithMaxSize(100*logrotate.MB),
			logrotate.WithRotateEvery(24*time.Hour),
			logrotate.WithMaxBackups(30),
			logrotate.WithMaxAge(30*logrotate.Day),
			logrotate.WithCompress(),
		)
		if err != nil {
			return nil, err
		}
		writers = append(writers, w)
		closer = w.Close
	}
	if config.Log.Output == "stdout" || config.Log.Output == "both" {
		writers = append(writers, os.Stdout)
	}

	log := slog.New(slog.NewJSONHandler(io.MultiWriter(writers...), &slog.HandlerOptions{
		Level: config.Log.SlogLevel(),
	}))
	slog.SetDefault(log)

	return &Logger{Logger: log, close: closer}, nil
}

// NewSlog unwraps the plain *slog.Logger for the rest of the app.
func NewSlog(i do.Injector) (*slog.Logger, error) {
	return do.MustInvoke[*Logger](i).Logger, nil
}
