package main

import (
	"GoSyntaxDoc/infrastructure/consumers"
	"GoSyntaxDoc/infrastructure/database"
	"GoSyntaxDoc/infrastructure/database/migrations"
	"GoSyntaxDoc/infrastructure/redis"
	"GoSyntaxDoc/infrastructure/repositories"
	"GoSyntaxDoc/presentation/middleware"
	"GoSyntaxDoc/services/user"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const maxRetries = 5                  // ✅ Maximum retry attempts for Kafka connection
const retryInterval = 5 * time.Second // ✅ Initial retry interval

func main() {
	fmt.Printf("Waiting 10 seconds for services start")
	time.Sleep(10 * time.Second)
	// ✅ Initialize Database
	database.ConnectDB()
	defer database.CloseDB()
	migrations.InitDB()

	// ✅ Initialize Redis
	redisService := redis.NewRedisService()
	defer redisService.Close() // ✅ Ensure Redis connection is closed

	// ✅ Initialize User Repository & Service
	userRepo := repositories.NewUserRepository(&database.Database)
	userService := user.NewUserService(userRepo)

	// ✅ Start Kafka Consumer with Retry Mechanism
	var kafkaConsumer *consumers.KafkaConsumer
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		kafkaConsumer, err = consumers.NewKafkaConsumer(
			[]string{"kafka:9092"},
			"user-service-group",
			[]string{"user.created", "user.fetch", "user.read"},
			userService,
			redisService,
		)
		if err == nil {
			fmt.Printf("✅ Kafka Consumer connected successfully on attempt %d\n", attempt)
			break
		}

		fmt.Printf("⚠ Failed to connect Kafka Consumer (attempt %d/%d): %v\n", attempt, maxRetries, err)
		if attempt == maxRetries {
			middleware.Log.Error("❌ Kafka Consumer failed after max retries, exiting.")
			os.Exit(1)
		}

		fmt.Printf("🔄 Retrying Kafka Consumer connection in %v...\n", retryInterval)
		time.Sleep(retryInterval)
	}

	// ✅ Run consumer in a separate goroutine
	go kafkaConsumer.ConsumeMessages()

	fmt.Println("🚀 Kafka Consumer running... Press Ctrl+C to stop.")

	// ✅ Graceful Shutdown Handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // Wait for shutdown signal

	// ✅ Close Kafka Consumer Gracefully
	fmt.Println("🛑 Shutting down Kafka Consumer microservice...")
	if err := kafkaConsumer.Close(); err != nil {
		middleware.Log.Error("Error closing Kafka consumer: ", err)
	}
}
