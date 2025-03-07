package gosyntaxdoc

import (
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const kafkaBroker = "kafka:9092"

var topics = []string{
	"notification.create",
	"notification.fetch",
	"notification.read",
}

// Function to check if a topic exists
func topicExists(topic string) bool {
	conn, err := kafka.Dial("tcp", kafkaBroker)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions(topic)
	if err != nil {
		return false // Topic does not exist
	}
	return len(partitions) > 0
}

// Function to create topics
func createKafkaTopics() {
	conn, err := kafka.Dial("tcp", kafkaBroker)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer conn.Close()

	for _, topic := range topics {
		if topicExists(topic) {
			log.Printf("✅ Topic '%s' already exists", topic)
			continue
		}

		err = conn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})

		if err != nil {
			log.Printf("⚠️ Failed to create topic '%s': %v", topic, err)
		} else {
			log.Printf("✅ Successfully created topic: %s", topic)
		}
	}
}

func main() {
	log.Println("⏳ Waiting for Kafka to be ready...")
	time.Sleep(10 * time.Second) // Wait for Kafka startup

	createKafkaTopics()
	log.Println("✅ Kafka topics setup completed!")
}
