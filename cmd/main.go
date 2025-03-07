package main

import (
	"GoSyntaxDoc/infrastructure"
	"GoSyntaxDoc/infrastructure/database"
	"GoSyntaxDoc/infrastructure/database/migrations"
	"GoSyntaxDoc/presentation/middleware"
	"GoSyntaxDoc/presentation/user"
	"GoSyntaxDoc/presentation/websocket"

	"github.com/gofiber/fiber/v2"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to the Greetings API!")
}

func main() {

	database.ConnectDB()
	defer database.CloseDB()
	producer := infrastructure.NewKafkaProducer([]string{"kafka:9092"}, "events")
	defer producer.Close()
	wsManager := websocket.NewWebSocketManager(producer)
	migrations.InitDB()
	app := fiber.New()
	app.Use(middleware.FiberLogger())
	app.Use(middleware.RecoveryMiddleware())
	app.Get("/", welcome)
	user.RegisterUserRoutes(app)
	websocket.RegisterWebsocketRoutes(app, wsManager)
	app.Get("/fatal", func(c *fiber.Ctx) error {
		panic("This is a fatal error")
	})

	app.Listen(":3002")
}
