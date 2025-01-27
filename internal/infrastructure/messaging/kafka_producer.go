package messaging

import (
	"ChallengeUALA/internal/platform/config"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
)

// KafkaProducer es un productor de eventos de Kafka.
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer crea una nueva instancia de KafkaProducer.
func NewKafkaProducer(kafkaConfig config.KafkaConfig) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(kafkaConfig.Brokers...),
			Topic:    kafkaConfig.Topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// PublishEvent publica un evento en un tópico de Kafka.
func (kp *KafkaProducer) PublishEvent(ctx context.Context, key string, value any) error {

	serializedValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshalling value: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: serializedValue,
	}

	return kp.writer.WriteMessages(ctx, msg)
}
