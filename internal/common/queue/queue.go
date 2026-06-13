package queue

import (
	"context"
	"encoding/json"

	"github.com/yash492/statusy/internal/common/jsonutil"
)

// Message represents a single raw message from the queue.
type Message struct {
	ID      string          `json:"id"`
	Payload json.RawMessage `json:"payload"` // Holds the raw JSON []byte from the broker
}

// MessageEnvelope wraps a message with a type-safe payload.
type MessageEnvelope[T any] struct {
	ID      string
	Payload T
}

// UnmarshalMessage decodes the raw JSON bytes of Message.Payload into a type-safe MessageEnvelope[T].
func UnmarshalMessage[T any](msg Message) (MessageEnvelope[T], error) {
	payload, err := jsonutil.UnmarshalWithType[T](msg.Payload)
	if err != nil {
		return MessageEnvelope[T]{}, err
	}
	return MessageEnvelope[T]{
		ID:      msg.ID,
		Payload: payload,
	}, nil
}

type EventType string

const (
	EventTypeIncidentUpdate    EventType = "incident_update"
	EventTypeMaintenanceUpdate EventType = "maintenance_update"
)

// AlertPayload represents the payload enqueued for notification dispatch.
type AlertPayload struct {
	EventType EventType `json:"event_type"`
	EventID   uint      `json:"event_id"`
}

// Queue defines the interface for interacting with the message queue.
type Queue interface {
	// Send serializes the payload to JSON and sends it to the specified queue.
	// Returns the message ID of the enqueued message.
	Send(ctx context.Context, queueName string, payload json.RawMessage) (string, error)

	// SendBatch serializes multiple payloads to JSON and sends them as a batch.
	// Returns the message IDs of the enqueued messages.
	SendBatch(ctx context.Context, queueName string, payloads []json.RawMessage) ([]string, error)

	// Read reads up to 'limit' messages with a visibility timeout of 'vt' seconds.
	Read(ctx context.Context, queueName string, vt int, limit int) ([]Message, error)

	// Delete removes a message from the queue after successful processing.
	Delete(ctx context.Context, queueName string, messageID string) error

	// Archive moves a message to the archive/dead-letter queue.
	Archive(ctx context.Context, queueName string, messageID string) error
}
