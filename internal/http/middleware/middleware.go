package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
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
		recover.New(),
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
