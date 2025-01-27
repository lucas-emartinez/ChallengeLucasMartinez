package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ChallengeUALA/internal/application/services"
	"ChallengeUALA/internal/domain"
	"ChallengeUALA/internal/infrastructure/messaging"
	"ChallengeUALA/internal/infrastructure/repositories"
	"ChallengeUALA/internal/platform/config"

	"github.com/go-redis/redis/v8"
)

func main() {
	log.Println("Starting consumer...")
	
	cfg, err := config.LoadAppConfig()
	if err != nil {
		log.Fatalf("Error loading app config: %v", err)
	}

	// Redis
	redisClient := redis.NewClient(cfg.Redis)
	redisRepo := repositories.NewRedisRepository(redisClient)

	timelineService := services.NewTimelineService(redisRepo)

	// Kafka
	kafka := messaging.NewKafkaConsumer(cfg.Kafka)
	defer kafka.Reader.Close()

	// Canal para manejar señales de terminación (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx := context.Background()

	// Goroutine para escuchar eventos de Kafka
	go func() {
		for {
			msg, err := kafka.Reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message from Kafka: %v", err)
				return
			}

			var tweet domain.Tweet
			if err := json.Unmarshal(msg.Value, &tweet); err != nil {
				log.Printf("Error unmarshalling tweet: %v", err)
				continue
			}

			// Aca actualizo el timeline de los seguidores del usuario que publico el tweet
			if err := timelineService.UpdateTimeline(ctx, &tweet); err != nil {
				log.Printf("Error updating timeline: %v", err)
			} else {
				log.Printf("Timeline updated for followers of user %s", tweet.UserID)
			}
		}
	}()

	<-sigChan
	log.Println("Shutting down consumer...")
}
