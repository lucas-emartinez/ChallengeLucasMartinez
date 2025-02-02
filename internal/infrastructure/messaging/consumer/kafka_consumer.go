package consumer

import (
	"ChallengeUALA/internal/platform/config"
	"github.com/segmentio/kafka-go"
)

// KafkaConsumer es un consumidor de eventos de Kafka.
type KafkaConsumer struct {
	Reader *kafka.Reader
}

// NewKafkaConsumer crea una nueva instancia de KafkaConsumer.
func NewKafkaConsumer(kafkaConfig config.KafkaConfig) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaConfig.Brokers,
		Topic:   kafkaConfig.Topic,
	})

	return &KafkaConsumer{
		Reader: reader,
	}
}
