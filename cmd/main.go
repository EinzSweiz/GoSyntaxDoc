package main

import (
	"GoSyntaxDoc/infrastructure"
	"GoSyntaxDoc/infrastructure/redis"
	"GoSyntaxDoc/presentation/middleware"
	"GoSyntaxDoc/presentation/websocket"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// ✅ Initialize Database
	// database.ConnectDB()
	// defer database.CloseDB()
	// migrations.InitDB()

	// ✅ Initialize Kafka Producer (for event-driven communication)
	producer := infrastructure.NewKafkaProducer([]string{"kafka:9092"})
	redisService := redis.NewRedisService()
	// ✅ Initialize User Repository & Service
	// userRepo := repositories.NewUserRepository(&database.Database)
	// userService := user.NewUserService(userRepo)

	// ✅ Start Kafka Consumer (Event-driven processing)
	// kafkaConsumer, err := consumers.NewKafkaConsumer(
	// 	[]string{"kafka:9092"},                              // Kafka brokers
	// 	"user-service-group",                                // Kafka Consumer Group ID
	// 	[]string{"user.created", "user.fetch", "user.read"}, // Kafka topics
	// 	userService,
	// 	redisService,
	// )
	// if err != nil {
	// 	fmt.Println("�� Failed to start Kafka consumer:", err)
	// 	os.Exit(1)
	// }

	// go kafkaConsumer.ConsumeMessages()

	// ✅ Set Up WebSocket Manager
	wsManager := websocket.NewWebSocketManager(producer, redisService)

	// ✅ Initialize Fiber (only for WebSockets)
	app := fiber.New()
	app.Use(middleware.FiberLogger())
	app.Use(middleware.RecoveryMiddleware())

	// ✅ Register WebSocket Routes
	websocket.RegisterWebsocketRoutes(app, wsManager)

	// ✅ Debug Route
	app.Get("/fatal", func(c *fiber.Ctx) error {
		panic("This is a fatal error")
	})

	// ✅ Start Fiber WebSocket Server in a Goroutine
	go func() {
		fmt.Println("🚀 WebSocket server is running on :3002")
		if err := app.Listen(":3002"); err != nil {
			fmt.Println("❌ Error starting WebSocket server:", err)
		}
	}()

	// ✅ Graceful Shutdown Handling
	fmt.Println("🚀 Event-driven microservice running... Press Ctrl+C to stop.")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // Wait for shutdown signal

	fmt.Println("🛑 Shutting down microservice...")
}
