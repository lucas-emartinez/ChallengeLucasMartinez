package messaging

import (
	"ChallengeUALA/internal/platform/config"
	"context"
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
func (kp *KafkaProducer) PublishEvent(ctx context.Context, key string, value []byte) error {

	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}

	return kp.writer.WriteMessages(ctx, msg)
}
