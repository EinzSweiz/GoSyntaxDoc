package infrastructure

import (
	"GoSyntaxDoc/presentation/middleware"
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type KafkaProducer struct {
	Brokers []string
}

// ✅ Constructor Function (Dependency Injection)
func NewKafkaProducer(brokers []string) *KafkaProducer {
	return &KafkaProducer{
		Brokers: brokers,
	}
}

// ✅ Send Message to Dynamic Topics
func (p *KafkaProducer) ProduceMessage(topic string, key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ✅ Create a new writer for each topic dynamically
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:   p.Brokers,
		Topic:     topic, // ✅ Dynamic Topic
		Balancer:  &kafka.LeastBytes{},
		BatchSize: 1,
	})

	message := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	err := writer.WriteMessages(ctx, message)
	if err != nil {
		middleware.Log.WithFields(logrus.Fields{
			"error": err,
			"topic": topic,
		}).Error("❌ Error producing Kafka message")
		return err
	}

	middleware.Log.WithFields(logrus.Fields{
		"message": "✅ Kafka Message Sent",
		"topic":   topic,
		"value":   value,
	}).Info("✅ Message sent successfully to Kafka")

	// ✅ Close writer after sending the message
	return writer.Close()
}
