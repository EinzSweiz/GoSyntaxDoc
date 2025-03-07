package websocket

import (
	"GoSyntaxDoc/infrastructure"
	"GoSyntaxDoc/presentation/middleware"
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/sirupsen/logrus"
)

type WebsocketEvent struct {
	Type  string      `json:"type"`
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type WebSocketManager struct {
	Producer *infrastructure.KafkaProducer
}

func NewWebSocketManager(producer *infrastructure.KafkaProducer) *WebSocketManager {
	return &WebSocketManager{Producer: producer}
}

func (wsm *WebSocketManager) HandleWebsocket(c *websocket.Conn) {
	defer c.Close()

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
		err = wsm.Producer.ProduceMessage(event.Event, string(message))
		if err != nil {
			middleware.Log.WithFields(logrus.Fields{"error": err}).Error("Error producing Kafka message")
			continue
		}
		middleware.Log.WithFields(logrus.Fields{"event": event.Event}).Info("WebSocket message received and forwarded to Kafka")
	}
}
