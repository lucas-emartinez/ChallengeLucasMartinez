package producer

import (
	"ChallengeUALA/internal/platform/config"
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKafkaWriter struct {
	mock.Mock
}

func (m *MockKafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	args := m.Called(ctx, msgs)
	return args.Error(0)
}

// Close simula el cierre del escritor.
func (m *MockKafkaWriter) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestKafkaProducer_PublishEvent(t *testing.T) {

	mockWriter := new(MockKafkaWriter)
	kp := &KafkaProducer{
		KafkaWriter: mockWriter,
	}

	mockWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(nil)

	err := kp.PublishEvent(context.Background(), "test-key", []byte("test-value"))

	assert.NoError(t, err)

	mockWriter.AssertExpectations(t)
}

func TestKafkaProducer_PublishEvent_Error(t *testing.T) {

	mockWriter := new(MockKafkaWriter)
	kp := &KafkaProducer{
		KafkaWriter: mockWriter,
	}

	mockWriter.On("WriteMessages", mock.Anything, mock.Anything).Return(assert.AnError)

	err := kp.PublishEvent(context.Background(), "test-key", []byte("test-value"))

	assert.Error(t, err)

	mockWriter.AssertExpectations(t)
}

func TestNewKafkaProducer(t *testing.T) {
	kafkaConfig := config.KafkaConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test-topic",
	}

	kp := NewKafkaProducer(kafkaConfig)

	assert.NotNil(t, kp)
}

func TestKafkaProducer_Writer(t *testing.T) {
	kafkaConfig := config.KafkaConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test-topic",
	}

	kp := NewKafkaProducer(kafkaConfig)

	writer := kp.Writer()

	assert.NotNil(t, writer)
}
