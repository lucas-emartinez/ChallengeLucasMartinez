package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"ChallengeUALA/internal/application/ports"
	"ChallengeUALA/internal/domain"
)

type TimelineService struct {
	tweetRepo  ports.TweetRepository
	followRepo ports.FollowRepository
	redisRepo  ports.RedisRepository
	dlq        ports.DeadLetterQueue
	logger     *log.Logger
}

func NewTimelineService(
	tweetRepo ports.TweetRepository,
	followRepo ports.FollowRepository,
	redisRepo ports.RedisRepository,
	dlq ports.DeadLetterQueue,
	logger *log.Logger,
) *TimelineService {
	return &TimelineService{
		tweetRepo:  tweetRepo,
		followRepo: followRepo,
		redisRepo:  redisRepo,
		dlq:        dlq,
		logger:     logger,
	}
}

func (s *TimelineService) UpdateTimeline(ctx context.Context, tweet *domain.Tweet) error {

	followers, err := s.followRepo.GetFollowers(ctx, tweet.UserID)
	if err != nil {
		payload, err := json.Marshal(tweet)
		if err != nil {
			return fmt.Errorf("error getting followers: %w", err)
		}

		_ = s.dlq.StoreEvent(ctx, "timeline_events", payload)

		return fmt.Errorf("error getting followers: %w", err)
	}

	for _, followerID := range followers {
		if err := s.redisRepo.AddToTimeline(ctx, followerID, tweet); err != nil {
			payload, err := json.Marshal(tweet)
			_ = s.dlq.StoreEvent(ctx, "timeline_events", payload)
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
		return nil, fmt.Errorf("error in calling redisRepo.GetTimeLine: %w", err)
	}

	return tweets, nil
}
