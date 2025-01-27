package ports

import (
	"ChallengeUALA/internal/domain"
	"context"
)

// TweetRepository define el contrato para almacenar y recuperar tweets (puerto de salida)
type TweetRepository interface {
	Save(ctx context.Context, tweet *domain.Tweet) error
}

type FollowRepository interface {
	Follow(ctx context.Context, followerID string, followedID string) error
	IsFollowing(ctx context.Context, followerID string, followedID string) (bool, error)
	GetFollowers(ctx context.Context, userID string) ([]string, error)
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, error)
}

type RedisRepository interface {
	AddToTimeline(ctx context.Context, userID string, tweet *domain.Tweet) error
	GetTimeline(ctx context.Context, userID string) ([]*domain.Tweet, error)
}
