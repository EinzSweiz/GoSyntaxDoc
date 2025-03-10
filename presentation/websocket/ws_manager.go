package websocket

import (
	"GoSyntaxDoc/infrastructure"
	"GoSyntaxDoc/infrastructure/redis"
	"GoSyntaxDoc/presentation/middleware"
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/sirupsen/logrus"
)

type WebsocketEvent struct {
	Type  string      `json:"type"`
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type WebSocketManager struct {
	Producer     *infrastructure.KafkaProducer
	RedisService *redis.RedisService
	clients      map[*websocket.Conn]bool
	mu           sync.Mutex
}

func NewWebSocketManager(producer *infrastructure.KafkaProducer, redisService *redis.RedisService) *WebSocketManager {
	wsm := &WebSocketManager{
		Producer:     producer,
		RedisService: redisService,
		clients:      make(map[*websocket.Conn]bool),
	}
	go wsm.listenToRedis()
	return wsm
}

func (wsm *WebSocketManager) HandleWebsocket(c *websocket.Conn) {
	defer c.Close()

	wsm.mu.Lock()
	wsm.clients[c] = true
	wsm.mu.Unlock()
	logrus.Infof("✅ WebSocket client registered: %v", c.RemoteAddr())

	defer func() {
		wsm.mu.Lock()
		delete(wsm.clients, c)
		wsm.mu.Unlock()
	}()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			middleware.Log.WithFields(logrus.Fields{"error": err}).Error("Error reading message from WebSocket")
			break
		}

		var event WebsocketEvent
		err = json.Unmarshal(message, &event)
		if err != nil {
			middleware.Log.WithFields(logrus.Fields{"error": err}).Error("Error unmarshalling JSON message")
			continue
		}

		// ✅ Construct the Kafka topic dynamically
		kafkaTopic := event.Event + "." + event.Type

		// ✅ Send to Kafka with correct topic
		err = wsm.Producer.ProduceMessage(kafkaTopic, event.Event, string(message))
		if err != nil {
			middleware.Log.WithFields(logrus.Fields{"error": err}).Error("Error producing Kafka message")
			continue
		}
		middleware.Log.WithFields(logrus.Fields{"topic": kafkaTopic}).Info("✅ WebSocket message forwarded to Kafka")
	}
}

func (wsm *WebSocketManager) listenToRedis() {
	wsm.RedisService.Subscribe("users_actions", func(msg string) {
		logrus.Infof("✅ Received Redis message: %s", msg) // ✅ Debug log

		wsm.mu.Lock()
		defer wsm.mu.Unlock()

		for client := range wsm.clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				logrus.WithFields(logrus.Fields{"error": err}).Error("❌ Error writing message to WebSocket")
				delete(wsm.clients, client) // Remove disconnected clients
			} else {
				logrus.Infof("✅ Message sent to WebSocket client: %v", client.RemoteAddr())
			}
		}
	})
}
