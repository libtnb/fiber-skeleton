package route

import (
	"github.com/gofiber/fiber/v3"
)

type Ws struct{}

func NewWs() *Ws {
	return &Ws{}
}

func (r *Ws) Register(router fiber.Router) {
	// github.com/gofiber/contrib/websocket middleware not support fiber v3 now
	/*router.Get("/ws", websocket.New(func(c *websocket.Conn) {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				return
			}
			if err = c.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		}
	}))*/
}
