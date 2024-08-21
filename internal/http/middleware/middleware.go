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

	"github.com/TheTNB/go-web-skeleton/internal/app"
)

// GlobalMiddleware is a collection of global middleware that will be applied to every request.
func GlobalMiddleware() []fiber.Handler {
	return []fiber.Handler{
		recover.New(),
		cors.New(),
		compress.New(),
		etag.New(),
		helmet.New(),
		requestid.New(),
		logger.New(),
		encryptcookie.New(encryptcookie.Config{
			Key: app.Conf.String("app.key"),
		}),
	}
}
