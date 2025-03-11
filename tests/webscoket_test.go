package websocket_test

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"GoSyntaxDoc/infrastructure"
	"GoSyntaxDoc/infrastructure/redis"
	wsm "GoSyntaxDoc/presentation/websocket"
)

// ✅ Mock Kafka Producer (Only If You Want to Avoid Real Kafka)
type MockKafkaProducer struct{}

func (m *MockKafkaProducer) ProduceMessage(topic, key, value string) error {
	return nil
}

func TestWebSocketWithRealRedisAndKafka(t *testing.T) {
	// ✅ Start Real Kafka (Ensure it's running at localhost:9092)
	realKafkaProducer := &infrastructure.KafkaProducer{
		Brokers: []string{"localhost:9092"},
	}

	// ✅ Initialize Fiber app
	app := fiber.New()

	// ✅ Use Real RedisService (make sure Redis is running!)
	realRedis := redis.NewRedisService()
	defer realRedis.Close() // Cleanup Redis connection

	// ✅ Create WebSocketManager with real Kafka and Redis
	manager := wsm.NewWebSocketManager(realKafkaProducer, realRedis)

	// ✅ Register WebSocket route
	wsm.RegisterWebsocketRoutes(app, manager)

	// ✅ Start test server
	go func() {
		_ = app.Listen(":8080")
	}()
	defer app.Shutdown()

	// ✅ Allow time for server to start
	time.Sleep(1 * time.Second)

	// ✅ Use Gorilla WebSocket Client
	client, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	require.NoError(t, err)
	defer client.Close()

	// ✅ Send test message
	message := `{"type": "created", "event": "user", "data": {"first_name": "RAID", "last_name": "Suline"}}`
	err = client.WriteMessage(websocket.TextMessage, []byte(message))
	require.NoError(t, err)

	// ✅ Publish message to Redis (simulating external event)
	err = realRedis.Publish("users_actions", `{"type": "created", "event": "user", "data": {"first_name": "RAID", "last_name": "Suline"}}`)
	require.NoError(t, err)

	// ✅ Read response (if any)
	_, response, err := client.ReadMessage()
	if err == nil {
		assert.NotEmpty(t, response)
	}
}
