package services_test

import (
	"ChallengeUALA/internal/application/services"
	"ChallengeUALA/internal/domain"
	"ChallengeUALA/internal/infrastructure/dlq"
	"context"
	"errors"
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTweetRepository simula el almacenamiento de tweets.
type MockTweetRepository struct {
	mock.Mock
}

func (m *MockTweetRepository) Save(ctx context.Context, tweet *domain.Tweet) error {
	args := m.Called(ctx, tweet)
	return args.Error(0)
}

// MockEventProducer simula la publicación de eventos en Kafka.
type MockEventProducer struct {
	mock.Mock
	wg *sync.WaitGroup // Para sincronizar la goroutine
}

func (m *MockEventProducer) PublishEvent(ctx context.Context, userID string, payload []byte) error {
	defer m.wg.Done() // Reduce el contador cuando se ejecuta
	args := m.Called(ctx, userID, payload)
	return args.Error(0)
}

// MockDeadLetterQueue simula una cola de eventos fallidos.
type MockDeadLetterQueue struct {
	mock.Mock
	wg *sync.WaitGroup // Para sincronizar la goroutine
}

func (m *MockDeadLetterQueue) StoreEvent(ctx context.Context, eventType string, payload []byte) error {
	defer m.wg.Done() // Reduce el contador cuando se ejecuta
	args := m.Called(ctx, eventType, payload)
	return args.Error(0)
}

func (m *MockDeadLetterQueue) RetrieveEvents(ctx context.Context) ([]dlq.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dlq.Message), args.Error(1)
}

func TestPostTweet_Success(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	content := "Hello, world!"
	logger := log.Default()

	mockRepo := new(MockTweetRepository)
	mockProducer := new(MockEventProducer)
	mockDLQ := new(MockDeadLetterQueue)

	var wg sync.WaitGroup
	wg.Add(1) // Se espera que al menos una goroutine termine

	mockProducer.wg = &wg
	mockDLQ.wg = &wg

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	mockProducer.On("PublishEvent", mock.Anything, userID, mock.Anything).Return(nil)

	tweetService := services.NewTweetService(mockRepo, mockProducer, mockDLQ, logger)

	err := tweetService.PostTweet(ctx, userID, content)
	assert.NoError(t, err)

	wg.Wait()

	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockDLQ.AssertNotCalled(t, "StoreEvent") // DLQ no debería ser llamado en éxito
}

func TestPostTweet_EventProducerFails(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	content := "Hello, world!"
	logger := log.Default()

	mockRepo := new(MockTweetRepository)
	mockProducer := new(MockEventProducer)
	mockDLQ := new(MockDeadLetterQueue)

	var wg sync.WaitGroup
	wg.Add(4) // Se espera 3 intentos + 1 guardado en DLQ

	mockProducer.wg = &wg
	mockDLQ.wg = &wg

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	mockProducer.On("PublishEvent", mock.Anything, userID, mock.Anything).Return(errors.New("publish failed")).Times(3)
	mockDLQ.On("StoreEvent", mock.Anything, "tweet_events", mock.Anything).Return(nil)

	tweetService := services.NewTweetService(mockRepo, mockProducer, mockDLQ, logger)

	err := tweetService.PostTweet(ctx, userID, content)
	assert.NoError(t, err)

	wg.Wait()

	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockDLQ.AssertExpectations(t)
}

func TestPostTweet_JsonMarshalFail(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	content := "Hello, world!"
	logger := log.Default()

	mockRepo := new(MockTweetRepository)
	mockProducer := new(MockEventProducer)
	mockDLQ := new(MockDeadLetterQueue)

	var wg sync.WaitGroup
	wg.Add(1) // Se espera 1 intento

	mockProducer.wg = &wg
	mockDLQ.wg = &wg

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	mockProducer.On("PublishEvent", mock.Anything, userID, mock.Anything).Return(nil)

	tweetService := services.NewTweetService(mockRepo, mockProducer, mockDLQ, logger)

	err := tweetService.PostTweet(ctx, userID, content)
	assert.NoError(t, err)

	wg.Wait()

	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockDLQ.AssertNotCalled(t, "StoreEvent") // DLQ no debería ser llamado en éxito
}

func TestPostTweet_StoreEventFails(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	content := "Hello, world!"
	logger := log.Default()

	mockRepo := new(MockTweetRepository)
	mockProducer := new(MockEventProducer)
	mockDLQ := new(MockDeadLetterQueue)

	var wg sync.WaitGroup
	wg.Add(4) // 3 intentos de PublishEvent + 1 intento de StoreEvent

	mockProducer.wg = &wg
	mockDLQ.wg = &wg

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)
	mockProducer.On("PublishEvent", mock.Anything, userID, mock.Anything).Return(errors.New("publish failed")).Times(3)
	mockDLQ.On("StoreEvent", mock.Anything, "tweet_events", mock.Anything).Return(errors.New("store event failed"))

	tweetService := services.NewTweetService(mockRepo, mockProducer, mockDLQ, logger)

	err := tweetService.PostTweet(ctx, userID, content)
	assert.NoError(t, err)

	wg.Wait()

	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockDLQ.AssertExpectations(t)
}

func TestPostTweet_NewTweetFails(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	invalidContent := ""
	logger := log.Default()

	mockRepo := new(MockTweetRepository)
	mockProducer := new(MockEventProducer)
	mockDLQ := new(MockDeadLetterQueue)

	tweetService := services.NewTweetService(mockRepo, mockProducer, mockDLQ, logger)

	err := tweetService.PostTweet(ctx, userID, invalidContent)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tweet content is empty")

	mockRepo.AssertNotCalled(t, "Save")
	mockProducer.AssertNotCalled(t, "PublishEvent")
	mockDLQ.AssertNotCalled(t, "StoreEvent")
}

func TestPostTweet_NewTeetExceedsLength(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	invalidContent := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam auctor, massa nec efficitur gravida, or " +
		"ci magna tincidunt nunc, nec ultricies nunc nunc nec nunc. Donec as sfsf ssfsf sfsf sf sfsf sf s fsfs s s s fsfsfs " +
		"ultr alk klakak ka kak aka kfjalfalfalfl laflalfallalfaalfalslfalflaslfsalflasflaslflasflkaskjfsakfjaskfa"

	logger := log.Default()

	mockRepo := new(MockTweetRepository)
	mockProducer := new(MockEventProducer)
	mockDLQ := new(MockDeadLetterQueue)

	tweetService := services.NewTweetService(mockRepo, mockProducer, mockDLQ, logger)

	err := tweetService.PostTweet(ctx, userID, invalidContent)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "tweet content is too long")

	mockRepo.AssertNotCalled(t, "Save")
	mockProducer.AssertNotCalled(t, "PublishEvent")
	mockDLQ.AssertNotCalled(t, "StoreEvent")
}

func TestPostTweet_SaveFails(t *testing.T) {
	ctx := context.Background()
	userID := "user123"
	content := "Hello, world!"
	logger := log.Default()

	mockRepo := new(MockTweetRepository)
	mockProducer := new(MockEventProducer)
	mockDLQ := new(MockDeadLetterQueue)

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(errors.New("db error"))

	tweetService := services.NewTweetService(mockRepo, mockProducer, mockDLQ, logger)

	err := tweetService.PostTweet(ctx, userID, content)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error saving tweet")

	mockProducer.AssertNotCalled(t, "PublishEvent")
	mockDLQ.AssertNotCalled(t, "StoreEvent")
}
