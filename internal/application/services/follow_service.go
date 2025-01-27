package services

import (
	"ChallengeUALA/internal/application/ports"
	"context"
	"errors"
	"fmt"
)

type FollowService struct {
	followRepo ports.FollowRepository
	userRepo   ports.UserRepository
}

func NewFollowService(followRepo ports.FollowRepository, userRepo ports.UserRepository) *FollowService {
	return &FollowService{
		followRepo: followRepo,
		userRepo:   userRepo,
	}
}

// Follow permite a un usuario seguir a otro
func (s *FollowService) Follow(ctx context.Context, followerID, followeeID string) error {

	if followerID == followeeID {
		return errors.New("a user cannot follow themselves")
	}

	// Seguidor
	_, err := s.userRepo.GetByID(ctx, followerID)
	if err != nil {
		return fmt.Errorf("error fetching follower: %w", err)
	}

	// Al que se quiere seguir
	_, err = s.userRepo.GetByID(ctx, followeeID)
	if err != nil {
		return fmt.Errorf("error fetching followee: %w", err)
	}

	// Acá validamos si el usuario ya sigue al otro
	isFollowing, err := s.followRepo.IsFollowing(ctx, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("error calling IsFollowing(): %w", err)
	}
	if isFollowing {
		return fmt.Errorf("user %s is already following user %s", followerID, followeeID)
	}

	if err := s.followRepo.Follow(ctx, followerID, followeeID); err != nil {
		return fmt.Errorf("error following user: %w", err)
	}

	return nil
}
