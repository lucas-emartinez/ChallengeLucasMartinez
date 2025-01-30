package main

import (
	"ChallengeUALA/internal/application/services"
	"ChallengeUALA/internal/domain"
	"ChallengeUALA/internal/infrastructure/dlq"
	"ChallengeUALA/internal/infrastructure/messaging/consumer"
	"ChallengeUALA/internal/infrastructure/messaging/producer"
	"ChallengeUALA/internal/infrastructure/repositories"
	"ChallengeUALA/internal/interfaces/http"
	"ChallengeUALA/internal/platform/config"
	"ChallengeUALA/internal/worker"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Cargar la configuración de la aplicación
	cfg, err := config.LoadAppConfig()
	if err != nil {
		log.Fatalf("Error loading app config: %v", err)
	}

	// Crear logger
	logger := log.New(os.Stdout, "App: ", log.LstdFlags|log.Lshortfile)

	// Configuración de Redis
	redisClient := redis.NewClient(cfg.Redis)
	redisRepo := repositories.NewRedisRepository(redisClient)

	// Inicializar los repositorios
	userRepository := repositories.NewUserRepository()
	tweetRepository := repositories.NewTweetRepository()
	followRepository := repositories.NewFollowRepository()

	// KafkaProducer
	kafkaProducer := producer.NewKafkaProducer(cfg.Kafka)

	// DLQ
	deadLetterQueue := dlq.NewDLQ()

	// Servicios
	tweetService := services.NewTweetService(tweetRepository, kafkaProducer, deadLetterQueue, logger)
	followService := services.NewFollowService(followRepository, userRepository)
	timelineService := services.NewTimelineService(tweetRepository, followRepository, redisRepo, deadLetterQueue, logger)

	// Configuración de Fiber para la API
	app := fiber.New()

	// Setup de las rutas de la API
	http.SetupRoutes(app, tweetService, followService, timelineService)

	// Iniciar la API en una goroutine
	go func() {
		log.Fatal(app.Listen(":8080"))
	}()

	// Kafka Consumer
	kafkaConsumer := consumer.NewKafkaConsumer(cfg.Kafka)

	// Canal para manejar señales de terminación
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine para manejar mensajes de Kafka
	go func() {
		log.Println("Kafka consumer started")
		ctx := context.Background()
		for {
			msg, err := kafkaConsumer.Reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}
			log.Println("Message received: ", string(msg.Value))
			var tweet domain.Tweet
			if err := json.Unmarshal(msg.Value, &tweet); err != nil {
				log.Printf("Error unmarshalling tweet: %v", err)
				continue
			}

			// Actualizar el timeline para los seguidores del usuario que publicó el tweet
			if err := timelineService.UpdateTimeline(ctx, &tweet); err != nil {
				log.Printf("Error updating timeline: %v", err)
			} else {
				log.Printf("Timeline updated for followers of user %s", tweet.UserID)
			}
		}
	}()

	// Dead Letter Queue Worker
	dlqWorker := worker.NewDLQWorker(deadLetterQueue, producer.NewKafkaProducer(cfg.Kafka), logger)
	go dlqWorker.Start(context.Background())

	<-sigChan
	log.Println("Shutting down...")
}
