package consumer_test

import (
	"ChallengeUALA/internal/infrastructure/messaging/consumer"
	"ChallengeUALA/internal/platform/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKafkaConsumer(t *testing.T) {

	kafkaConfig := config.KafkaConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test-topic",
	}

	kafkaConsumer := consumer.NewKafkaConsumer(kafkaConfig)

	assert.NotNil(t, kafkaConsumer)
}
