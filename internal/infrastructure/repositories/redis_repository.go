package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"ChallengeUALA/internal/domain"
	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

// AddToTimeline agrega un tweet al timeline de un usuario usando LPUSH.
func (r *RedisRepository) AddToTimeline(ctx context.Context, userID string, tweet *domain.Tweet) error {
	// Serializar el tweet a JSON
	log.Println("Adding to Timeline: ", userID)
	tweetJSON, err := json.Marshal(tweet)
	if err != nil {
		return fmt.Errorf("error marshalling tweet: %w", err)
	}

	// Agrega el tweet al inicio de la lista del timeline del usuario
	// ordenado por la fecha del tweet
	err = r.client.ZAdd(ctx, "timeline:"+userID, &redis.Z{
		Score:  float64(tweet.CreatedAt.Unix()),
		Member: tweetJSON,
	}).Err()
	if err != nil {
		return fmt.Errorf("error adding tweet to timeline: %w", err)
	}

	// Limitar el tamaño de la lista para evitar que crezca indefinidamente
	// Por ejemplo, mantener solo los últimos 1000 tweets
	err = r.client.ZRemRangeByRank(ctx, "timeline:"+userID, 0, -1001).Err()
	if err != nil {
		return fmt.Errorf("error trimming timeline: %w", err)
	}

	return nil
}

func (r *RedisRepository) GetTimeline(ctx context.Context, userID string) ([]*domain.Tweet, error) {
	// le devuelvo los últimos 100 tweets del timeline
	fmt.Println("Getting Timeline: ", userID)
	tweetsJSON, err := r.client.ZRevRange(ctx, "timeline:"+userID, 0, 99).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting timeline: %w", err)
	}

	if len(tweetsJSON) == 0 {
		return nil, fmt.Errorf("timeline is empty for user %s", userID)
	}

	var tweets []*domain.Tweet
	for _, tweetJSON := range tweetsJSON {
		var tweet domain.Tweet
		if err := json.Unmarshal([]byte(tweetJSON), &tweet); err != nil {
			log.Printf("Error unmarshalling tweet: %v", err)
			continue // si tenemos un error le vamos a mostrar el proximo tweet de igual manera.
		}
		tweets = append(tweets, &tweet)
	}

	log.Printf("Retrieved %d tweets for user %s", len(tweets), userID)
	return tweets, nil
}
