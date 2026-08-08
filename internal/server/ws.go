package server

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"

	"github.com/libtnb/fiber-skeleton/internal/pkg/transport"
)

func WsRoutes() transport.Endpoints {
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
	}
}
