package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/encryptcookie"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/google/wire"
	"github.com/knadh/koanf/v2"
)

var ProviderSet = wire.NewSet(NewMiddlewares)

type Middlewares struct {
	conf *koanf.Koanf
}

func NewMiddlewares(conf *koanf.Koanf) *Middlewares {
	return &Middlewares{
		conf: conf,
	}
}

// Globals is a collection of global middleware that will be applied to every request.
func (r *Middlewares) Globals(app *fiber.App) []fiber.Handler {
	return []fiber.Handler{
		recover.New(recover.Config{
			EnableStackTrace: true,
		}),
		cors.New(),
		compress.New(),
		etag.New(),
		helmet.New(),
		requestid.New(),
		logger.New(),
		encryptcookie.New(encryptcookie.Config{
			Key: r.conf.String("app.key"),
		}),
	}
}
