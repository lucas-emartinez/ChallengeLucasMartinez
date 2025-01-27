package main

import (
	"ChallengeUALA/internal/interfaces/http"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"log"

	"ChallengeUALA/internal/application/services"
	"ChallengeUALA/internal/infrastructure/dlq"
	"ChallengeUALA/internal/infrastructure/messaging"
	"ChallengeUALA/internal/infrastructure/repositories"
	"ChallengeUALA/internal/platform/config"
)

func main() {
	cfg, err := config.LoadAppConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	logger := log.New(log.Writer(), "API: ", log.LstdFlags)

	redisClient := redis.NewClient(cfg.Redis)
	redisRepo := repositories.NewRedisRepository(redisClient)

	eventProducer := messaging.NewKafkaProducer(cfg.Kafka)

	deadLetterQueue := dlq.NewDLQ()

	// Inicializar repositorios
	inMemoryUserRepository := repositories.NewUserRepository()
	inMemoryTweetRepository := repositories.NewTweetRepository()
	inMemoryFollowRepository := repositories.NewFollowRepository()

	// Servicios
	tweetService := services.NewTweetService(inMemoryTweetRepository, eventProducer, deadLetterQueue, logger)
	followService := services.NewFollowService(inMemoryFollowRepository, inMemoryUserRepository)
	timelineService := services.NewTimelineService(redisRepo)

	// app de fiber
	app := fiber.New()

	// Setup de las rutas
	http.SetupRoutes(app, tweetService, followService, timelineService)

	log.Fatal(app.Listen(":8080"))
}
