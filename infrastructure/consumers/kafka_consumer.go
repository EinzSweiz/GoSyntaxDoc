package consumers

import (
	"GoSyntaxDoc/domain/events"
	"GoSyntaxDoc/infrastructure/redis"
	"GoSyntaxDoc/services/user"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

const RedisChannel = "users_actions"

type KafkaConsumer struct {
	Reader       *kafka.Reader
	UserService  *user.UserService
	RedisService *redis.RedisService
}

// ✅ NewKafkaConsumer: Handles connection retries and proper initialization
func NewKafkaConsumer(brokers []string, groupID string, topics []string, userService *user.UserService, redisService *redis.RedisService) (*KafkaConsumer, error) {
	maxRetries := 5

	// ✅ Kafka Reader Configuration
	readerConfig := kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		GroupTopics:    topics,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        1 * time.Second,
		StartOffset:    kafka.FirstOffset,
		CommitInterval: time.Second, // ✅ Automatically commit offsets
	}

	// ✅ Retry Connecting to Kafka
	for attempt := 1; attempt <= maxRetries; attempt++ {
		dialer := &kafka.Dialer{
			Timeout:   5 * time.Second,
			DualStack: true,
		}

		conn, err := dialer.DialContext(context.Background(), "tcp", brokers[0])
		if err == nil {
			conn.Close()
			logrus.Infof("✅ Successfully connected to Kafka broker on attempt %d", attempt)
			return &KafkaConsumer{
				Reader:       kafka.NewReader(readerConfig),
				UserService:  userService,
				RedisService: redisService,
			}, nil
		}

		logrus.Warnf("⚠ Kafka connection attempt %d/%d failed: %v", attempt, maxRetries, err)

		if attempt == maxRetries {
			logrus.Error("❌ Maximum retry attempts reached, could not connect to Kafka")
			return nil, fmt.Errorf("failed to connect to Kafka after %d attempts", maxRetries)
		}

		time.Sleep(10 * time.Second) // 🔄 Wait before retrying
	}

	return nil, fmt.Errorf("unexpected error creating Kafka consumer")
}

// ✅ ConsumeMessages: Listens to Kafka and processes events
func (c *KafkaConsumer) ConsumeMessages() {
	if c == nil || c.Reader == nil {
		logrus.Error("❌ Kafka consumer is not properly initialized")
		return
	}

	logrus.Info("🚀 Kafka Consumer started and listening for messages...")

	for {
		msg, err := c.Reader.ReadMessage(context.Background())
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err}).Error("❌ Error reading Kafka message")
			time.Sleep(10 * time.Second) // Prevent log spam on failure
			continue
		}

		logrus.Infof("\n📩 Kafka Message Received:\nTopic: %s\nValue: %s\n", msg.Topic, string(msg.Value))

		switch msg.Topic {
		case "user.created":
			c.handleUserCreate(msg.Value)
		case "user.fetch":
			c.handleUserFetchById(msg.Value)
		case "user.read":
			c.handlerUserFetchAll(msg.Value)
		default:
			logrus.Infof("⚠️ Unsupported Kafka message topic: %s", msg.Topic)
			continue

		}

		// ✅ Commit the message to prevent reprocessing
		if err := c.Reader.CommitMessages(context.Background(), msg); err != nil {
			logrus.WithFields(logrus.Fields{"error": err}).Error("❌ Failed to commit message")
		}
	}
}

func (c *KafkaConsumer) handleUserCreate(value []byte) {
	var event events.KafkaUserCreatedEvent
	err := json.Unmarshal(value, &event)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("❌ Error unmarshalling JSON message")
		return
	}

	// ✅ Extract User Data
	firstName := event.Data.FirstName
	lastName := event.Data.LastName
	logrus.Infof("Extracted Data: FirstName=%s, LastName=%s", firstName, lastName)

	// ✅ Validate Data
	if firstName == "" || lastName == "" {
		logrus.Error("❌ Invalid user data in Kafka message")
		return
	}

	// ✅ Call Service Layer to Process Business Logic
	user, err := c.UserService.HandleUserCreated(firstName, lastName)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("❌ Failed to create user from Kafka event")
		return
	}

	c.publishToRedis("user.created", user)
}

func (c *KafkaConsumer) handleUserFetchById(value []byte) {
	var event events.KafkaUserFetchByIdEvent

	err := json.Unmarshal(value, &event)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("❌ Error unmarshalling JSON message")
		return
	}

	userId := event.Data.UserID
	logrus.Infof("Extracted Data: UserID=%d", userId)
	if userId <= 0 {
		logrus.Error("❌ Invalid user ID in Kafka message")
		return
	}

	// ✅ Call the service with a clean integer (not raw JSON)
	user, err := c.UserService.FetchUserById(userId)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("❌ Failed to fetch user from Kafka event")
		return
	}

	c.publishToRedis("user.fetch", user)
}

func (c *KafkaConsumer) handlerUserFetchAll(value []byte) {
	var event events.KafkaUserReadAllEvent

	err := json.Unmarshal(value, &event)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("�� Error unmarshalling JSON message")
		return
	}
	users, err := c.UserService.HandleUserRead()
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Error("�� Failed to fetch all users from Kafka event")
		return
	}

	c.publishToRedis("user.read", users)
}

// ✅ Publish to Redis (Reusable function)
func (c *KafkaConsumer) publishToRedis(event string, data interface{}) {
	// ✅ Convert data to JSON
	userData, err := json.Marshal(data)
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err}).Errorf("❌ Failed to marshal data for event: %s", event)
		return
	}

	// ✅ Publish asynchronously to Redis
	go func() {
		err := c.RedisService.Publish(RedisChannel, string(userData))
		if err != nil {
			logrus.WithFields(logrus.Fields{"error": err}).Errorf("❌ Failed to publish data to Redis for event: %s", event)
		} else {
			logrus.Infof("✅ Successfully published event [%s] to Redis: %s", event, userData)
		}
	}()
}

// ✅ Close: Gracefully shuts down the Kafka consumer
func (c *KafkaConsumer) Close() error {
	if c != nil && c.Reader != nil {
		return c.Reader.Close()
	}
	return nil
}
