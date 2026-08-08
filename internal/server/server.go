package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/libtnb/validator"
	"github.com/libtnb/validator/contrib/openapi"

	"github.com/libtnb/fiber-skeleton/internal/conf"
	"github.com/libtnb/fiber-skeleton/internal/pkg/registry"
)

func NewRouter(
	config *conf.Config,
	log *slog.Logger,
	validate *validator.Validator,
	version Version,
	routes registry.Routes,
) (*fiber.App, error) {
	r := fiber.New(fiber.Config{
		AppName:           config.App.Name,
		BodyLimit:         config.HTTP.BodyLimit << 10,
		ReadBufferSize:    config.HTTP.HeaderLimit,
		ReadTimeout:       config.HTTP.ReadTimeout,
		WriteTimeout:      config.HTTP.WriteTimeout,
		IdleTimeout:       config.HTTP.IdleTimeout,
		ReduceMemoryUsage: config.HTTP.ReduceMemoryUsage,
		// every framework-level error (404, 405, 413, panics) leaves as JSON
		ErrorHandler: errorHandler,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	for _, handler := range globalMiddlewares(config, log) {
		r.Use(handler)
	}

	HTTP(routes, r)

	if config.HTTP.Docs {
		spec, err := SpecJSON(config.App.Name, version, validate, routes)
		if err != nil {
			return nil, err
		}
		docs := openapi.DocsHTML(config.App.Name, "/openapi.json")
		r.Get("/openapi.json", func(c fiber.Ctx) error {
			c.Type("json")
			return c.Send(spec)
		})
		r.Get("/docs", func(c fiber.Ctx) error {
			c.Type("html")
			return c.Send(docs)
		})
	}

	return r, nil
}

// errorHandler is the single error exit; 5xx details are logged, not sent.
func errorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := errors.AsType[*fiber.Error](err); ok {
		code = e.Code
	}

	if code >= fiber.StatusInternalServerError {
		slog.ErrorContext(c.Context(), "unhandled error",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Any("err", err),
		)
		return c.Status(code).JSON(fiber.Map{"msg": http.StatusText(code)})
	}

	return c.Status(code).JSON(fiber.Map{"msg": err.Error()})
}
