package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGMQQueue struct {
	db *pgxpool.Pool
}

func NewPGMQQueue(db *pgxpool.Pool) *PGMQQueue {
	return &PGMQQueue{db: db}
}

var _ Queue = &PGMQQueue{}

func (q *PGMQQueue) Send(ctx context.Context, queueName string, payload json.RawMessage) (string, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	var msgID int64
	err = q.db.QueryRow(ctx, "SELECT pgmq.send($1, $2::jsonb);", queueName, string(bytes)).Scan(&msgID)
	if err != nil {
		return "", fmt.Errorf("failed to send pgmq message: %w", err)
	}
	return strconv.FormatInt(msgID, 10), nil
}

func (q *PGMQQueue) SendBatch(ctx context.Context, queueName string, payloads []json.RawMessage) ([]string, error) {
	if len(payloads) == 0 {
		return nil, nil
	}

	jsonStrings := make([]string, len(payloads))
	for i, p := range payloads {
		bytes, err := json.Marshal(p)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal batch payload at index %d: %w", i, err)
		}
		jsonStrings[i] = string(bytes)
	}

	rows, err := q.db.Query(ctx, "SELECT * FROM pgmq.send_batch($1, $2::jsonb[]);", queueName, jsonStrings)
	if err != nil {
		return nil, fmt.Errorf("failed to send pgmq batch: %w", err)
	}
	defer rows.Close()

	var msgIDs []string
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan sent msg id: %w", err)
		}
		msgIDs = append(msgIDs, strconv.FormatInt(id, 10))
	}
	return msgIDs, nil
}

func (q *PGMQQueue) Read(ctx context.Context, queueName string, vt int, limit int) ([]Message, error) {
	rows, err := q.db.Query(ctx, "SELECT msg_id, message FROM pgmq.read($1, $2, $3);", queueName, vt, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to read messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msgID int64
		var rawPayload []byte
		if err := rows.Scan(&msgID, &rawPayload); err != nil {
			return nil, fmt.Errorf("failed to scan queue message: %w", err)
		}
		messages = append(messages, Message{
			ID:      strconv.FormatInt(msgID, 10),
			Payload: rawPayload,
		})
	}
	return messages, nil
}

func (q *PGMQQueue) Delete(ctx context.Context, queueName string, messageID string) error {
	msgID, err := strconv.ParseInt(messageID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid message ID format: %w", err)
	}

	var deleted bool
	err = q.db.QueryRow(ctx, "SELECT pgmq.delete($1::text, $2::bigint);", queueName, msgID).Scan(&deleted)
	if err != nil {
		return fmt.Errorf("failed to delete message %d: %w", msgID, err)
	}
	return nil
}

func (q *PGMQQueue) Archive(ctx context.Context, queueName string, messageID string) error {
	msgID, err := strconv.ParseInt(messageID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid message ID format: %w", err)
	}

	var archived bool
	err = q.db.QueryRow(ctx, "SELECT pgmq.archive($1, $2);", queueName, msgID).Scan(&archived)
	if err != nil {
		return fmt.Errorf("failed to archive message %d: %w", msgID, err)
	}
	return nil
}
