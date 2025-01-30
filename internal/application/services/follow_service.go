package services

import (
	"ChallengeUALA/internal/application/ports"
	"context"
	"fmt"
	uuid2 "github.com/google/uuid"
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

	err := validateUUID(followerID, followeeID)
	if err != nil {
		return err
	}

	if followerID == followeeID {
		return fmt.Errorf("user %s can't follow itself", followerID)
	}

	// Seguidor
	_, err = s.userRepo.GetByID(ctx, followerID)
	if err != nil {
		return fmt.Errorf("error in calling userRepo.GetByID(): %w", err)
	}

	// Al que se quiere seguir
	_, err = s.userRepo.GetByID(ctx, followeeID)
	if err != nil {
		return fmt.Errorf("error in calling userRepo.GetByID(): %w", err)
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
		return fmt.Errorf("error in calling followRepo.Follow(): %w", err)
	}

	return nil
}

func validateUUID(uuid ...string) error {
	for _, u := range uuid {
		_, err := uuid2.Parse(u)
		if err != nil {
			return fmt.Errorf("invalid UUID: %s", u)
		}
	}
	return nil
}
