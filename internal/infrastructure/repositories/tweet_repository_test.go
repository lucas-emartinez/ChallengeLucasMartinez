package repositories

import (
	"context"
	"testing"

	"ChallengeUALA/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewTweetRepository(t *testing.T) {
	repo := NewTweetRepository()
	assert.NotNil(t, repo)
	assert.NotNil(t, repo.tweets)
}

func TestTweetRepository_Save(t *testing.T) {
	repo := NewTweetRepository()
	tweet := &domain.Tweet{ID: "1", Content: "Hello, world!"}

	err := repo.Save(context.Background(), tweet)
	assert.NoError(t, err)

	repo.mu.RLock()
	defer repo.mu.RUnlock()
	savedTweet, exists := repo.tweets[tweet.ID]
	assert.True(t, exists)
	assert.Equal(t, tweet, savedTweet)
}
