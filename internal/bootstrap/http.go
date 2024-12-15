package bootstrap

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/knadh/koanf/v2"

	"github.com/go-rat/fiber-skeleton/internal/http/middleware"
	"github.com/go-rat/fiber-skeleton/internal/route"
)

func NewRouter(conf *koanf.Koanf, middlewares *middleware.Middlewares, http *route.Http, ws *route.Ws) *fiber.App {
	// prefork not support dual stack
	network := fiber.NetworkTCP
	if conf.Bool("http.prefork") {
		network = fiber.NetworkTCP4
	}

	r := fiber.New(fiber.Config{
		AppName:               conf.String("app.name"),
		BodyLimit:             conf.MustInt("http.bodyLimit") << 10,
		ReadBufferSize:        conf.MustInt("http.headerLimit"),
		ReduceMemoryUsage:     conf.Bool("http.reduceMemoryUsage"),
		Prefork:               conf.Bool("http.prefork"),
		EnablePrintRoutes:     conf.Bool("http.debug"),
		DisableStartupMessage: !conf.Bool("http.debug"),
		Network:               network,
		// replace default json encoder and decoder if you are not happy with the performance
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// add middleware
	for _, handler := range middlewares.Globals(r) {
		r.Use(handler)
	}

	// add http route
	http.Register(r)
	// add ws route
	ws.Register(r)

	// add fallback handler
	r.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("404 Not Found")
	})

	return r
}
