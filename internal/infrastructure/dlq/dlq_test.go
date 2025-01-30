package dlq

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDLQ(t *testing.T) {
	// Crear una nueva DLQ
	dlq := NewDLQ()

	// Verificar que la instancia de DLQ no sea nil
	assert.NotNil(t, dlq)

	// Verificar que la longitud de los mensajes sea 0 al principio
	assert.Len(t, dlq.messages, 0)
}

func TestStoreEvent(t *testing.T) {
	// Crear una nueva DLQ
	dlq := NewDLQ()

	// Crear un evento
	eventType := "event-1"
	payload := []byte("This is a test event")

	err := dlq.StoreEvent(context.Background(), eventType, payload)

	assert.NoError(t, err)

	assert.Len(t, dlq.messages, 1)

	assert.Equal(t, eventType, dlq.messages[0].ID)
	assert.Equal(t, "This is a test event", string(dlq.messages[0].Body))
	assert.WithinDuration(t, time.Now(), dlq.messages[0].Timestamp, time.Second)
}

func TestStoreEventWhenQueueIsFull(t *testing.T) {
	dlq := NewDLQ()

	// LLeno la cola con 100 mensajes
	for i := 0; i < 100; i++ {
		eventType := "event-" + strconv.Itoa(i)
		payload := []byte("This is test event #" + strconv.Itoa(i))
		dlq.StoreEvent(context.Background(), eventType, payload)
	}

	assert.Len(t, dlq.messages, 100)

	eventType := "event-101"
	payload := []byte("This is test event #101")
	err := dlq.StoreEvent(context.Background(), eventType, payload)

	assert.NoError(t, err)

	assert.Len(t, dlq.messages, 100)
	assert.Equal(t, "event-1", dlq.messages[0].ID) // El mensaje más antiguo debería haber sido eliminado
}

func TestRetrieveEvents(t *testing.T) {
	dlq := NewDLQ()

	// Guardar eventos
	err := dlq.StoreEvent(context.Background(), "event-1", []byte("Event 1"))
	assert.NoError(t, err)
	// Recuperar eventos
	events, err := dlq.RetrieveEvents(context.Background())
	assert.NoError(t, err)
	assert.Len(t, events, 1)
	assert.Equal(t, string(events[0].Body), "Event 1")
}
