package infrastructure

import (
	"GoSyntaxDoc/presentation/middleware"
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

// ✅ Constructor Function (Dependency Injection)
func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:   brokers,
			Topic:     topic,
			Balancer:  &kafka.LeastBytes{},
			BatchSize: 1,
		}),
	}
}

func (p *KafkaProducer) ProduceMessage(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	err := p.Writer.WriteMessages(ctx, message)
	if err != nil {
		middleware.Log.WithFields(logrus.Fields{
			"error": err,
		})
		return err

	}
	middleware.Log.WithFields(logrus.Fields{
		"message": "✅ Kafka Message Sent",
		"value":   value,
	})
	return nil
}

func (p *KafkaProducer) Close() error {
	if p.Writer != nil {
		middleware.Log.Info("Kafka Producer closed")
		return p.Writer.Close()
	}
	return nil
}
