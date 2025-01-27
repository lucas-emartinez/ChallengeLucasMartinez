package messaging

import (
	"ChallengeUALA/internal/platform/config"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

// KafkaConsumer es un consumidor de eventos de Kafka.
type KafkaConsumer struct {
	Reader *kafka.Reader
}

// NewKafkaConsumer crea una nueva instancia de KafkaConsumer.
func NewKafkaConsumer(kafkaConfig config.KafkaConfig) (*KafkaConsumer, error) {
	log.Println("KafkaConfig: ", kafkaConfig)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kafkaConfig.Brokers,
		Topic:   kafkaConfig.Topic,
	})

	if reader == nil {
		return nil, fmt.Errorf("error creating Kafka reader")
	}

	return &KafkaConsumer{
		Reader: reader,
	}, nil
}
