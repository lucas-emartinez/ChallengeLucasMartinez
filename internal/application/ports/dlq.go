package ports

import (
	"ChallengeUALA/internal/infrastructure/dlq"
	"context"
)

// DeadLetterQueue define el contrato para guardar eventos fallidos (puerto de salida)
type DeadLetterQueue interface {
	StoreEvent(ctx context.Context, eventType string, payload []byte) error
	RetrieveEvents(ctx context.Context) ([]dlq.Message, error)
}
