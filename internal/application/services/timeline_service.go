package services

import (
	"context"
	"fmt"
	"log"

	"ChallengeUALA/internal/application/ports"
	"ChallengeUALA/internal/domain"
)

type TimelineService struct {
	tweetRepo  ports.TweetRepository
	followRepo ports.FollowRepository
	redisRepo  ports.RedisRepository
}

func NewTimelineService(
	tweetRepo ports.TweetRepository,
	followRepo ports.FollowRepository,
	redisRepo ports.RedisRepository,
) *TimelineService {
	return &TimelineService{
		tweetRepo:  tweetRepo,
		followRepo: followRepo,
		redisRepo:  redisRepo,
	}
}

func (s *TimelineService) UpdateTimeline(ctx context.Context, tweet *domain.Tweet) error {
	log.Println("Getting followers")
	log.Println("Tweet: ", tweet)
	log.Println("Tweet UserID: ", tweet.UserID)
	log.Println(s.followRepo)
	followers, err := s.followRepo.GetFollowers(ctx, tweet.UserID)
	log.Println("Followers: ", followers)
	if err != nil {
		return fmt.Errorf("error getting followers: %w", err)
	}

	for _, followerID := range followers {
		if err := s.redisRepo.AddToTimeline(ctx, followerID, tweet); err != nil {
			return fmt.Errorf("error adding tweet to timeline: %w", err)
		}
	}

	return nil
}

// GetTimeline obtiene los últimos tweets del timeline de un usuario.
func (s *TimelineService) GetTimeline(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID is required")
	}

	tweets, err := s.redisRepo.GetTimeline(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting timeline: %w", err)
	}

	return tweets, nil
}
