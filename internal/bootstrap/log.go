package bootstrap

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/libtnb/logrotate"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/config"
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
	conf := do.MustInvoke[*config.Config](i)

	var (
		writers []io.Writer
		closer  = func() error { return nil }
	)

	if conf.Log.Output == "file" || conf.Log.Output == "both" {
		w, err := logrotate.New(conf.Log.Path,
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
	if conf.Log.Output == "stdout" || conf.Log.Output == "both" {
		writers = append(writers, os.Stdout)
	}

	log := slog.New(slog.NewJSONHandler(io.MultiWriter(writers...), &slog.HandlerOptions{
		Level: conf.Log.SlogLevel(),
	}))
	slog.SetDefault(log)

	return &Logger{Logger: log, close: closer}, nil
}

// NewSlog unwraps the plain *slog.Logger for the rest of the app.
func NewSlog(i do.Injector) (*slog.Logger, error) {
	return do.MustInvoke[*Logger](i).Logger, nil
}
