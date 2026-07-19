package server

import (
	"encoding/base64"
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/encryptcookie"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/libtnb/fiber-skeleton/internal/conf"
)

func globalMiddlewares(config *conf.Config, log *slog.Logger) []fiber.Handler {
	handlers := []fiber.Handler{
		recover.New(recover.Config{
			EnableStackTrace: true,
		}),
	}

	// CORS only when origins are explicitly allowed; empty = same-origin
	if len(config.HTTP.CorsOrigins) > 0 {
		handlers = append(handlers, cors.New(cors.Config{
			AllowOrigins: config.HTTP.CorsOrigins,
		}))
	}

	return append(handlers,
		compress.New(),
		etag.New(),
		helmet.New(),
		requestid.New(),
		accessLog(log),
		encryptcookie.New(encryptcookie.Config{
			Key: base64.StdEncoding.EncodeToString([]byte(config.App.Key)),
		}),
	)
}

// accessLog writes one structured line per request through the app logger.
func accessLog(log *slog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		// probes are noise
		if c.Path() == "/healthz" || c.Path() == "/readyz" {
			return c.Next()
		}

		start := time.Now()
		err := c.Next()

		// the error handler runs later; reflect the status it will write
		status := c.Response().StatusCode()
		var fe *fiber.Error
		if errors.As(err, &fe) {
			status = fe.Code
		} else if err != nil {
			status = fiber.StatusInternalServerError
		}

		log.LogAttrs(c.Context(), slog.LevelInfo, "http request",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", status),
			slog.Duration("duration", time.Since(start)),
			slog.String("ip", c.IP()),
			slog.String("request_id", requestid.FromContext(c)),
		)

		return err
	}
}
