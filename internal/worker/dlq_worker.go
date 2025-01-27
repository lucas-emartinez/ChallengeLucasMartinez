package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"ChallengeUALA/internal/application/ports"
	"ChallengeUALA/internal/domain"
)

// DLQWorker es un worker que procesa mensajes de la DLQ.
type DLQWorker struct {
	deadLetterQueue ports.DeadLetterQueue
	eventProducer   ports.EventProducer
	logger          *log.Logger
}

// NewDLQWorker crea una nueva instancia de DLQWorker.
func NewDLQWorker(dlq ports.DeadLetterQueue, ep ports.EventProducer, logger *log.Logger) *DLQWorker {
	return &DLQWorker{
		deadLetterQueue: dlq,
		eventProducer:   ep,
		logger:          logger,
	}
}

// Start inicia el worker para procesar mensajes de la DLQ.
func (w *DLQWorker) Start(ctx context.Context) {
	// cada 5 minutos vamos a intentar reprocesar los mensajes de la DLQ
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Println("DLQ worker stopped")
			return
		case <-ticker.C:
			w.processDLQ(ctx)
		}
	}
}

// processDLQ recupera y reprocesa los mensajes de la DLQ.
func (w *DLQWorker) processDLQ(ctx context.Context) {
	messages, err := w.deadLetterQueue.RetrieveEvents(ctx)
	if err != nil {
		w.logger.Printf("Error retrieving events from DLQ: %v", err)
		return
	}

	for _, msg := range messages {
		var tweet domain.Tweet
		if err := json.Unmarshal([]byte(msg.Body), &tweet); err != nil {
			w.logger.Printf("Error unmarshalling tweet from DLQ: %v", err)
			continue
		}

		// Intentar publicar el evento nuevamente
		if err := w.eventProducer.PublishEvent(ctx, tweet.UserID, []byte(msg.Body)); err != nil {
			w.logger.Printf("Error reprocessing event from DLQ: %v", err)
			continue
		}

		w.logger.Printf("Successfully reprocessed event from DLQ: %s", msg.ID)
	}
}
