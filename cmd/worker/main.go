package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ChallengeUALA/internal/infrastructure/dlq"
	"ChallengeUALA/internal/infrastructure/messaging"
	"ChallengeUALA/internal/platform/config"
	"ChallengeUALA/internal/worker"
)

func main() {
	log.Println("Starting DLQ worker...")

	// Cargar la configuración de la aplicación
	cfg, err := config.LoadAppConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	// Inicializar el logger
	logger := log.New(os.Stdout, "DLQWorker: ", log.LstdFlags|log.Lshortfile)

	// Inicializar la DLQ
	dlq := dlq.NewDLQ()

	// Inicializar el productor de eventos de Kafka
	eventProducer := messaging.NewKafkaProducer(cfg.Kafka)

	// Crear una instancia del worker
	dlqWorker := worker.NewDLQWorker(dlq, eventProducer, logger)

	// Contexto para manejar la cancelación
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Canal para manejar señales de terminación (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar el worker en una goroutine
	go dlqWorker.Start(ctx)

	// Esperar una señal de terminación
	<-sigChan
	logger.Println("Shutting down DLQ worker...")
}
