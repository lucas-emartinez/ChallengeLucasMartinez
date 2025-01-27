package ports

import (
	"context"
)

// EventProducer define el contrato para enviar eventos (puerto de salida)
type EventProducer interface {
	PublishEvent(ctx context.Context, key string, value []byte) error
}
