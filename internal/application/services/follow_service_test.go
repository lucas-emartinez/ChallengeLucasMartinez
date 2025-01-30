package services_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ChallengeUALA/internal/application/services"
	"ChallengeUALA/internal/domain"
)

// Mock de FollowRepository
type MockFollowsRepository struct {
	mock.Mock
}

func (m *MockFollowsRepository) Follow(ctx context.Context, followerID, followedID string) error {
	args := m.Called(ctx, followerID, followedID)
	return args.Error(0)
}

func (m *MockFollowsRepository) IsFollowing(ctx context.Context, followerID, followedID string) (bool, error) {
	args := m.Called(ctx, followerID, followedID)
	return args.Bool(0), args.Error(1)
}

func (m *MockFollowsRepository) GetFollowers(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

// Mock de UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestFollowService_Follow(t *testing.T) {
	// Arrange
	mockFollowRepo := new(MockFollowsRepository)
	mockUserRepo := new(MockUserRepository)
	service := services.NewFollowService(mockFollowRepo, mockUserRepo)

	ctx := context.Background()
	followerID := uuid.NewString()
	followeeID := uuid.NewString()

	mockUserRepo.On("GetByID", ctx, followerID).Return(&domain.User{ID: followerID}, nil)
	mockUserRepo.On("GetByID", ctx, followeeID).Return(&domain.User{ID: followeeID}, nil)
	mockFollowRepo.On("IsFollowing", ctx, followerID, followeeID).Return(false, nil)
	mockFollowRepo.On("Follow", ctx, followerID, followeeID).Return(nil)

	// Act
	err := service.Follow(ctx, followerID, followeeID)

	// Assert
	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockFollowRepo.AssertExpectations(t)
}

func TestFollowService_Follow_SameUser(t *testing.T) {
	mockFollowRepo := new(MockFollowsRepository)
	mockUserRepo := new(MockUserRepository)
	service := services.NewFollowService(mockFollowRepo, mockUserRepo)

	ctx := context.Background()
	userID := uuid.NewString()

	err := service.Follow(ctx, userID, userID)
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("user %s can't follow itself", userID), err.Error())
}

func TestFollowService_Follow_AlreadyFollowing(t *testing.T) {
	// Arrange
	mockFollowRepo := new(MockFollowsRepository)
	mockUserRepo := new(MockUserRepository)
	service := services.NewFollowService(mockFollowRepo, mockUserRepo)

	ctx := context.Background()
	followerID := uuid.NewString()
	followeeID := uuid.NewString()

	mockUserRepo.On("GetByID", ctx, followerID).Return(&domain.User{ID: followerID}, nil)
	mockUserRepo.On("GetByID", ctx, followeeID).Return(&domain.User{ID: followeeID}, nil)
	mockFollowRepo.On("IsFollowing", ctx, followerID, followeeID).Return(true, nil)

	// Act
	err := service.Follow(ctx, followerID, followeeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("user %s is already following user %s", followerID, followeeID), err.Error())
}

func TestFollowService_Follow_InvalidUUID(t *testing.T) {
	// Arrange
	mockFollowRepo := new(MockFollowsRepository)
	mockUserRepo := new(MockUserRepository)
	service := services.NewFollowService(mockFollowRepo, mockUserRepo)

	ctx := context.Background()
	followerID := "asd-uuid"
	followeeID := uuid.NewString()

	// Act
	err := service.Follow(ctx, followerID, followeeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("invalid UUID: %s", followerID), err.Error())
}

func TestFollowService_Follow_ErrorFollowingUser(t *testing.T) {
	// Arrange
	mockFollowRepo := new(MockFollowsRepository)
	mockUserRepo := new(MockUserRepository)
	service := services.NewFollowService(mockFollowRepo, mockUserRepo)

	ctx := context.Background()
	followerID := uuid.NewString()
	followeeID := uuid.NewString()

	mockUserRepo.On("GetByID", ctx, followerID).Return(&domain.User{ID: followerID}, nil)
	mockUserRepo.On("GetByID", ctx, followeeID).Return(&domain.User{ID: followeeID}, nil)
	mockFollowRepo.On("IsFollowing", ctx, followerID, followeeID).Return(false, nil)
	mockFollowRepo.On("Follow", ctx, followerID, followeeID).Return(errors.New("error following user"))

	// Act
	err := service.Follow(ctx, followerID, followeeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "error in calling followRepo.Follow(): error following user", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestFollowService_Follow_GetFollowerError(t *testing.T) {
	// Arrange
	mockFollowRepo := new(MockFollowsRepository)
	mockUserRepo := new(MockUserRepository)
	service := services.NewFollowService(mockFollowRepo, mockUserRepo)

	ctx := context.Background()
	followerID := uuid.NewString()
	followeeID := uuid.NewString()

	mockUserRepo.On("GetByID", ctx, followerID).Return(&domain.User{}, errors.New("error getting user"))

	// Act
	err := service.Follow(ctx, followerID, followeeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "error in calling userRepo.GetByID(): error getting user", err.Error())
	mockUserRepo.AssertExpectations(t)
}

func TestFollowService_Follow_GetFolloweeError(t *testing.T) {
	// Arrange
	mockFollowRepo := new(MockFollowsRepository)
	mockUserRepo := new(MockUserRepository)
	service := services.NewFollowService(mockFollowRepo, mockUserRepo)

	ctx := context.Background()
	followerID := uuid.NewString()
	followeeID := uuid.NewString()

	mockUserRepo.On("GetByID", ctx, followerID).Return(&domain.User{ID: followerID}, nil)
	mockUserRepo.On("GetByID", ctx, followeeID).Return(&domain.User{}, errors.New("error getting user"))

	// Act
	err := service.Follow(ctx, followerID, followeeID)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "error in calling userRepo.GetByID(): error getting user", err.Error())
	mockUserRepo.AssertExpectations(t)
}
