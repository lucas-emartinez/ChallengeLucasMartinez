package producer

import (
	"ChallengeUALA/internal/platform/config"
	"context"

	"github.com/segmentio/kafka-go"
)

// KafkaWriter define una interfaz para mockear kafka.Writer
type KafkaWriter interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

// KafkaProducer es un productor de eventos de Kafka.
type KafkaProducer struct {
	KafkaWriter KafkaWriter
}

// NewKafkaProducer crea una nueva instancia de KafkaProducer.
func NewKafkaProducer(kafkaConfig config.KafkaConfig) *KafkaProducer {
	return &KafkaProducer{
		KafkaWriter: &kafka.Writer{
			Addr:     kafka.TCP(kafkaConfig.Brokers...),
			Topic:    kafkaConfig.Topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// PublishEvent publica un evento en Kafka.
func (kp *KafkaProducer) PublishEvent(ctx context.Context, key string, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
	}

	return kp.KafkaWriter.WriteMessages(ctx, msg)
}

// Writer devuelve el writer interno (solo para pruebas)
func (kp *KafkaProducer) Writer() KafkaWriter {
	return kp.KafkaWriter
}
