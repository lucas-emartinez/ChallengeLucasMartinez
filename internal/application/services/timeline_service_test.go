package services_test

import (
	"ChallengeUALA/internal/application/services"
	"ChallengeUALA/internal/infrastructure/dlq"
	"context"
	"errors"
	"testing"

	"ChallengeUALA/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockFollowRepository struct{ mock.Mock }
type MockRedisRepository struct{ mock.Mock }
type MockDLQ struct{ mock.Mock }

func (m *MockFollowRepository) Follow(ctx context.Context, followerID string, followedID string) error {
	args := m.Called(ctx, followerID, followedID)
	return args.Error(0)
}

func (m *MockFollowRepository) IsFollowing(ctx context.Context, followerID string, followedID string) (bool, error) {
	args := m.Called(ctx, followerID, followedID)
	return args.Bool(0), args.Error(1)
}

func (m *MockFollowRepository) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisRepository) AddToTimeline(ctx context.Context, userID string, tweet *domain.Tweet) error {
	args := m.Called(ctx, userID, tweet)
	return args.Error(0)
}

func (m *MockRedisRepository) GetTimeline(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.Tweet), args.Error(1)
}

func (m *MockDLQ) StoreEvent(ctx context.Context, queue string, payload []byte) error {
	args := m.Called(ctx, queue, payload)
	return args.Error(0)
}

func (m *MockDLQ) RetrieveEvents(ctx context.Context) ([]dlq.Message, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dlq.Message), args.Error(1)
}

// 🔹 Test UpdateTimeline (success)
func TestUpdateTimeline_Success(t *testing.T) {
	ctx := context.Background()
	mockFollowRepo := new(MockFollowRepository)
	mockRedisRepo := new(MockRedisRepository)
	mockDLQ := new(MockDLQ)

	tweet := &domain.Tweet{UserID: "user123", Content: "Hello world"}

	mockFollowRepo.On("GetFollowers", ctx, "user123").Return([]string{"follower1", "follower2"}, nil)
	mockRedisRepo.On("AddToTimeline", ctx, "follower1", tweet).Return(nil)
	mockRedisRepo.On("AddToTimeline", ctx, "follower2", tweet).Return(nil)

	err := services.NewTimelineService(nil, mockFollowRepo, mockRedisRepo, mockDLQ, nil).
		UpdateTimeline(ctx, tweet)

	assert.NoError(t, err)
	mockFollowRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}

// 🔹 Test UpdateTimeline - GetFollowers fails
func TestUpdateTimeline_GetFollowersFails(t *testing.T) {
	ctx := context.Background()
	mockFollowRepo := new(MockFollowRepository)
	mockRedisRepo := new(MockRedisRepository)
	mockDLQ := new(MockDLQ)

	service := services.NewTimelineService(nil, mockFollowRepo, mockRedisRepo, mockDLQ, nil)

	tweet := &domain.Tweet{UserID: "user123", Content: "Hello world"}
	mockFollowRepo.On("GetFollowers", ctx, "user123").Return([]string{}, errors.New("DB error"))
	mockDLQ.On("StoreEvent", ctx, "timeline_events", mock.Anything).Return(nil)

	err := service.UpdateTimeline(ctx, tweet)
	assert.Error(t, err)
	mockDLQ.AssertExpectations(t)
}

// 🔹 Test UpdateTimeline - AddToTimeline fails
func TestUpdateTimeline_AddToTimelineFails(t *testing.T) {
	ctx := context.Background()
	mockFollowRepo := new(MockFollowRepository)
	mockRedisRepo := new(MockRedisRepository)
	mockDLQ := new(MockDLQ)

	service := services.NewTimelineService(nil, mockFollowRepo, mockRedisRepo, mockDLQ, nil)

	tweet := &domain.Tweet{UserID: "user123", Content: "Hello world"}

	mockFollowRepo.On("GetFollowers", ctx, "user123").Return([]string{"follower1"}, nil)
	mockRedisRepo.On("AddToTimeline", ctx, "follower1", tweet).Return(errors.New("Redis error"))
	mockDLQ.On("StoreEvent", ctx, "timeline_events", mock.Anything).Return(nil)

	err := service.UpdateTimeline(ctx, tweet)
	assert.Error(t, err)
	mockDLQ.AssertExpectations(t)
}

// 🔹 Test GetTimeline - Success
func TestGetTimeline_Success(t *testing.T) {
	ctx := context.Background()
	mockRedisRepo := new(MockRedisRepository)
	service := services.NewTimelineService(nil, nil, mockRedisRepo, nil, nil)

	tweets := []*domain.Tweet{{UserID: "user123", Content: "Hello world"}}
	mockRedisRepo.On("GetTimeline", ctx, "user123").Return(tweets, nil)

	result, err := service.GetTimeline(ctx, "user123")
	assert.NoError(t, err)
	assert.Equal(t, tweets, result)

	mockRedisRepo.AssertExpectations(t)
}

// 🔹 Test GetTimeline - Redis failure
func TestGetTimeline_RedisFails(t *testing.T) {
	ctx := context.Background()
	mockRedisRepo := new(MockRedisRepository)
	service := services.NewTimelineService(nil, nil, mockRedisRepo, nil, nil)

	mockRedisRepo.On("GetTimeline", ctx, "user123").Return([]*domain.Tweet{}, errors.New("Redis error"))

	result, err := service.GetTimeline(ctx, "user123")
	assert.Error(t, err)
	assert.Nil(t, result)

	mockRedisRepo.AssertExpectations(t)
}

// 🔹 Test GetTimeline - UserID empty
func TestGetTimeline_UserIDEmpty(t *testing.T) {
	ctx := context.Background()
	service := services.NewTimelineService(nil, nil, nil, nil, nil)

	result, err := service.GetTimeline(ctx, "")
	assert.Error(t, err)
	assert.Nil(t, result)
}
