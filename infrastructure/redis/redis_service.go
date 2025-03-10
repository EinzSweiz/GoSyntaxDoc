package redis

import (
	"GoSyntaxDoc/presentation/middleware"
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisService() *RedisService {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		middleware.Log.WithFields(logrus.Fields{"error": err}).Error("Error connecting to Redis")
	} else {
		middleware.Log.Info("Connected to Redis")
	}
	return &RedisService{client: client, ctx: ctx}
}

func (rs *RedisService) Publish(channel string, message string) error {
	err := rs.client.Publish(rs.ctx, channel, message).Err()
	if err != nil {
		middleware.Log.WithFields(logrus.Fields{"error": err}).Error("Error publishing to Redis")
		return err
	}
	middleware.Log.Info("�� Published message to Redis channel: " + channel)
	return nil
}

func (r *RedisService) Subscribe(channel string, messageHandler func(msg string)) {
	pubsub := r.client.Subscribe(r.ctx, channel)

	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(r.ctx)
			if err != nil {
				middleware.Log.WithFields(logrus.Fields{"error": err}).Error("❌ Redis Subscription Error")

				// ✅ If Redis client is closed, reconnect automatically
				if err == redis.ErrClosed {
					middleware.Log.Warn("⚠️ Redis connection lost. Reconnecting...")
					pubsub = r.client.Subscribe(r.ctx, channel)
					continue
				}
				break // Exit on unexpected errors
			}

			// ✅ Process message asynchronously
			go messageHandler(msg.Payload)
		}
	}()
}

func (rs *RedisService) Close() {
	rs.client.Close()
	middleware.Log.Info("Redis connection closed")
}
