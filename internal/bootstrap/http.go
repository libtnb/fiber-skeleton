package bootstrap

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/libtnb/validator"
	"github.com/libtnb/validator/contrib/openapi"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/config"
	"github.com/libtnb/fiber-skeleton/internal/middleware"
	"github.com/libtnb/fiber-skeleton/internal/route"
)

func NewRouter(i do.Injector) (*fiber.App, error) {
	conf := do.MustInvoke[*config.Config](i)
	middlewares := do.MustInvoke[*middleware.Middlewares](i)

	// handlers reach this instance through service.Bind / validator.Default
	validator.SetDefault(do.MustInvoke[*validator.Validator](i))

	r := fiber.New(fiber.Config{
		AppName:           conf.App.Name,
		BodyLimit:         conf.HTTP.BodyLimit << 10,
		ReadBufferSize:    conf.HTTP.HeaderLimit,
		ReadTimeout:       conf.HTTP.ReadTimeout,
		WriteTimeout:      conf.HTTP.WriteTimeout,
		IdleTimeout:       conf.HTTP.IdleTimeout,
		ReduceMemoryUsage: conf.HTTP.ReduceMemoryUsage,
		// every framework-level error (404, 405, 413, panics) leaves as JSON
		ErrorHandler: errorHandler,
		// swap in a faster JSON codec here if it ever shows up in profiles
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	for _, handler := range middlewares.Globals(r) {
		r.Use(handler)
	}

	if err := route.HTTP(i, r); err != nil {
		return nil, err
	}

	if conf.HTTP.Docs {
		spec, err := route.SpecJSON(i, conf.App.Name)
		if err != nil {
			return nil, err
		}
		docs := openapi.DocsHTML(conf.App.Name, "/openapi.json")
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
	var e *fiber.Error
	if errors.As(err, &e) {
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
