package dlq

import (
	"context"
	"sync"
	"time"
)

// Message representa un mensaje que puede ser almacenado en la DLQ
type Message struct {
	ID        string
	Body      []byte
	Timestamp time.Time
	Error     error
}

// DLQ es una cola en memoria para almacenar mensajes fallidos
// En un entorno productivo, esto debería ser un servicio externo como SQS por ejemplo.
type DLQ struct {
	mu       sync.Mutex
	messages []Message
}

// NewDLQ crea una nueva instancia de DLQ
func NewDLQ() *DLQ {
	return &DLQ{
		messages: make([]Message, 0),
	}
}

// StoreEvent agrega un evento a la DLQ
func (d *DLQ) StoreEvent(ctx context.Context, eventType string, payload []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Si la cola está llena, elimino el mensaje más antiguo
	// Esto es solo para mantener la memoria por las dudas en un entorno de prueba
	if len(d.messages) >= 100 {
		d.messages = d.messages[1:]
	}

	// Uso el tipo de evento como el ID del mensaje
	msg := Message{
		ID:        eventType,
		Body:      payload,
		Timestamp: time.Now(),
	}

	d.messages = append(d.messages, msg)
	return nil
}

// RetrieveEvents obtiene todos los mensajes de la DLQ
func (d *DLQ) RetrieveEvents(ctx context.Context) ([]Message, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.messages, nil
}
