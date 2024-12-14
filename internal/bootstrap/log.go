package bootstrap

import (
	"log/slog"

	"github.com/knadh/koanf/v2"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewLog(conf *koanf.Koanf) *slog.Logger {
	ljLogger := &lumberjack.Logger{
		Filename: "storage/logs/app.log",
		MaxSize:  10,
		MaxAge:   30,
		Compress: true,
	}

	level := slog.LevelInfo
	if conf.Bool("app.debug") {
		level = slog.LevelDebug
	}

	log := slog.New(slog.NewJSONHandler(ljLogger, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(log)

	return log
}
