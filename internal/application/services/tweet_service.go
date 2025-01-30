package services

import (
	"ChallengeUALA/internal/application/ports"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"ChallengeUALA/internal/domain"
)

// TweetService es un servicio de aplicación que maneja la lógica de negocio relacionada con los tweets.
type TweetService struct {
	tweetRepo       ports.TweetRepository
	eventProducer   ports.EventProducer
	deadLetterQueue ports.DeadLetterQueue
	logger          *log.Logger
}

// NewTweetService crea una nueva instancia de TweetService
func NewTweetService(
	tr ports.TweetRepository,
	ep ports.EventProducer,
	dlq ports.DeadLetterQueue,
	logger *log.Logger,
) *TweetService {
	return &TweetService{
		tweetRepo:       tr,
		eventProducer:   ep,
		deadLetterQueue: dlq,
		logger:          logger,
	}
}

// PostTweet crea un nuevo tweet y lo guarda en la base de datos (En este caso, en memoria, pero sería POSTGRES)
// También envía un evento a Kafka para notificar que se creó un nuevo tweet.
func (s *TweetService) PostTweet(ctx context.Context, userID, content string) error {

	if len(content) > 280 {
		return fmt.Errorf("tweet content is too long")
	}

	if len(content) == 0 {
		return fmt.Errorf("tweet content is empty")
	}

	tweet := domain.NewTweet(userID, content)
	if err := s.tweetRepo.Save(ctx, tweet); err != nil {
		return fmt.Errorf("error saving tweet: %w", err)
	}

	go func() {
		// Nuevo contexto para el evento, con un timeout de 10 segundos.
		// Incluso si el cliente cancela la solicitud HTTP original con un timeout o aborta,
		// el evento se envia de todas formas.
		eventCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		const maxRetries = 3
		for i := 0; i < maxRetries; i++ {

			serializedTweet, _ := json.Marshal(tweet)
			err := s.eventProducer.PublishEvent(eventCtx, userID, serializedTweet)
			if err == nil {
				s.logger.Printf("Tweet event published")
				return
			}
			// Backoff exponencial, en este caso lo puse porque sirve como mecanismo de reintentos
			// sin embargo, no es la mejor práctica para manejar errores de red.
			// En un sistema real, se debería implementar un mecanismo de reintentos más robusto.
			// Por ejemplo, usando un circuit breaker o un retry con un backoff más inteligente.
			time.Sleep(time.Second * time.Duration(math.Pow(2, float64(i))))
		}

		// Si falla después de reintentos, guardar en DLQ
		// Un marshall a una struct no deberia fallar siempre y cuando la struct sea correcta
		payload, _ := json.Marshal(tweet)

		if err := s.deadLetterQueue.StoreEvent(ctx, "tweet_events", payload); err != nil {
			s.logger.Printf("Error storing event in DLQ: %v", err)
			return
		}
	}()

	return nil
}
