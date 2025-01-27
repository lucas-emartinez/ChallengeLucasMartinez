package repositories

import (
	"context"
	"sync"

	"ChallengeUALA/internal/domain"
)

// TweetRepository es una struct que implementa la interfaz TweetRepository
type TweetRepository struct {
	mu     sync.RWMutex
	tweets map[string]*domain.Tweet
}

// NewTweetRepository crea una nueva instancia de TweetRepository
func NewTweetRepository() *TweetRepository {
	return &TweetRepository{
		tweets: make(map[string]*domain.Tweet),
	}
}

// Save guarda un tweet en memoria
func (r *TweetRepository) Save(ctx context.Context, tweet *domain.Tweet) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tweets[tweet.ID] = tweet
	return nil
}
