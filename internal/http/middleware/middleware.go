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
	"github.com/knadh/koanf/v2"
)

// GlobalMiddleware is a collection of global middleware that will be applied to every request.
func GlobalMiddleware(conf *koanf.Koanf) []fiber.Handler {
	return []fiber.Handler{
		recover.New(),
		cors.New(),
		compress.New(),
		etag.New(),
		helmet.New(),
		requestid.New(),
		logger.New(),
		encryptcookie.New(encryptcookie.Config{
			Key: conf.String("app.key"),
		}),
	}
}
