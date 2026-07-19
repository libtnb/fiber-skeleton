package server

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/samber/do/v2"

	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
)

func WsRoutes(i do.Injector) (transport.Endpoints, error) {
	return transport.Endpoints{
		{Method: fiber.MethodGet, Path: "/ws", Handler: websocket.New(func(c *websocket.Conn) {
			for {
				_, msg, err := c.ReadMessage()
				if err != nil {
					return
				}
				if err = c.WriteMessage(websocket.TextMessage, msg); err != nil {
					return
				}
			}
		})},
	}, nil
}
