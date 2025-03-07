package websocket

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func RegisterWebsocketRoutes(app *fiber.App, manager *WebSocketManager) {
	app.Get("/ws", websocket.New(manager.HandleWebsocket))
}
