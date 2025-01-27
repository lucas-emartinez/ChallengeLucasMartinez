package config

import (
	"github.com/go-redis/redis/v8"
	"os"
)

type Config struct {
	Kafka KafkaConfig
	Redis *redis.Options
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

func LoadAppConfig() (*Config, error) {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	redisAddr := os.Getenv("REDIS_ADDR")

	kafkaConfig := KafkaConfig{
		Brokers: []string{kafkaBrokers},
		Topic:   kafkaTopic,
	}

	redisConfig := redis.Options{
		Addr: redisAddr,
		DB:   0,
	}

	return &Config{
		Kafka: kafkaConfig,
		Redis: &redisConfig,
	}, nil
}
