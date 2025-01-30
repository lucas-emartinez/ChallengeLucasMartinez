package repositories

import (
	"context"
	"log"
	"testing"
	"time"

	"ChallengeUALA/internal/domain"
	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

var pool *dockertest.Pool
var resource *dockertest.Resource

func setupTestRedisClient() (*redis.Client, func()) {
	var client *redis.Client
	var err error

	if pool == nil {
		pool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not connect to Docker: %s", err)
		}
	}

	if resource == nil {
		resource, err = pool.Run("redis", "latest", nil)
		if err != nil {
			log.Fatalf("Could not start resource: %s", err)
		}
	}

	if err := pool.Retry(func() error {
		client = redis.NewClient(&redis.Options{
			Addr: "localhost:" + resource.GetPort("6379/tcp"),
			DB:   1,
		})
		return client.Ping(context.Background()).Err()
	}); err != nil {
		log.Fatalf("Could not connect to Redis: %s", err)
	}

	// Return the client and a cleanup function
	return client, func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}

func TestRedisRepository_AddToTimeline(t *testing.T) {
	client, cleanup := setupTestRedisClient()
	defer cleanup()

	repo := NewRedisRepository(client)
	ctx := context.Background()

	tweet := &domain.Tweet{
		ID:        "1",
		Content:   "Hello, world!",
		CreatedAt: time.Now(),
	}

	err := repo.AddToTimeline(ctx, "user1", tweet)
	assert.NoError(t, err)

	// Verify that the tweet was added correctly
	tweets, err := repo.GetTimeline(ctx, "user1")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tweets))
	assert.Equal(t, tweet.Content, tweets[0].Content)
}

func TestRedisRepository_GetTimeline(t *testing.T) {
	client, cleanup := setupTestRedisClient()
	defer cleanup()

	repo := NewRedisRepository(client)
	ctx := context.Background()

	// Add some tweets to the timeline
	tweets := []*domain.Tweet{
		{ID: "1", Content: "Tweet 1", CreatedAt: time.Now().Add(-1 * time.Hour)},
		{ID: "2", Content: "Tweet 2", CreatedAt: time.Now()},
	}

	for _, tweet := range tweets {
		err := repo.AddToTimeline(ctx, "user2", tweet)
		assert.NoError(t, err)
	}

	// Retrieve the timeline and verify the tweets
	retrievedTweets, err := repo.GetTimeline(ctx, "user2")
	assert.NoError(t, err)
	assert.Equal(t, len(tweets), len(retrievedTweets))
	assert.Equal(t, tweets[0].Content, retrievedTweets[1].Content)
	assert.Equal(t, tweets[1].Content, retrievedTweets[0].Content)
}
