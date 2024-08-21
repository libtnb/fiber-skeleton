package bootstrap

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3"

	"github.com/TheTNB/go-web-skeleton/internal/app"
	"github.com/TheTNB/go-web-skeleton/internal/http/middleware"
	"github.com/TheTNB/go-web-skeleton/internal/route"
)

func initHttp() {
	app.Http = fiber.New(fiber.Config{
		AppName:           app.Conf.String("app.name"),
		BodyLimit:         app.Conf.MustInt("http.bodyLimit") << 10,
		ReadBufferSize:    app.Conf.MustInt("http.headerLimit"),
		ReduceMemoryUsage: app.Conf.Bool("http.reduceMemoryUsage"),
		// replace default json encoder and decoder if you are not happy with the performance
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// add middleware
	for _, handler := range middleware.GlobalMiddleware() {
		app.Http.Use(handler)
	}

	// add route
	route.Http(app.Http)

	// add fallback handler
	app.Http.Use(func(c fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("404 Not Found")
	})

	if err := app.Http.Listen(app.Conf.MustString("http.address"), fiber.ListenConfig{
		ListenerNetwork:       "tcp",
		EnablePrefork:         app.Conf.Bool("http.prefork"),
		EnablePrintRoutes:     app.Conf.Bool("http.debug"),
		DisableStartupMessage: !app.Conf.Bool("http.debug"),
	}); err != nil {
		panic(fmt.Sprintf("failed to start http server: %v", err))
	}
}
