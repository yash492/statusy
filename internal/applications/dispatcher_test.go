package applications

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/common/queue"
	"github.com/yash492/statusy/internal/domain/notifications"
)

// Mock notifications.NotificationsRepository implementation
type mockNotificationsRepo struct {
	mu           sync.Mutex
	deliveries   map[string]notifications.NotificationDelivery
	incidentCfgs []notifications.ViewNotification
	incDetails   notifications.IncidentNotificationDetails
	getDelErr    error
}

func (m *mockNotificationsRepo) Save(ctx context.Context, vn notifications.ViewNotification) (notifications.ViewNotification, error) {
	return vn, nil
}
func (m *mockNotificationsRepo) GetByViewID(ctx context.Context, viewID uint, search string, limit int, offset int) ([]notifications.ViewNotification, int64, error) {
	return nil, 0, nil
}
func (m *mockNotificationsRepo) Delete(ctx context.Context, id uint) error {
	return nil
}
func (m *mockNotificationsRepo) Update(ctx context.Context, vn notifications.ViewNotification) (notifications.ViewNotification, error) {
	return vn, nil
}
func (m *mockNotificationsRepo) GetDelivery(ctx context.Context, channelID uint, alertType string, alertID uint) (notifications.NotificationDelivery, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.getDelErr != nil {
		return notifications.NotificationDelivery{}, m.getDelErr
	}
	key := fmt.Sprintf("%d-%s-%d", channelID, alertType, alertID)
	d, ok := m.deliveries[key]
	if !ok {
		return notifications.NotificationDelivery{}, apperrors.NotFoundError("delivery not found", nil)
	}
	return d, nil
}
func (m *mockNotificationsRepo) SaveDelivery(ctx context.Context, delivery notifications.NotificationDelivery) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := fmt.Sprintf("%d-%s-%d", delivery.ViewNotificationID, delivery.AlertType, delivery.AlertID)
	m.deliveries[key] = delivery
	return nil
}
func (m *mockNotificationsRepo) UpdateDelivery(ctx context.Context, deliveryID uint, lastUpdateID uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, d := range m.deliveries {
		if d.ID == deliveryID {
			d.LastUpdateID = lastUpdateID
			m.deliveries[k] = d
			break
		}
	}
	return nil
}
func (m *mockNotificationsRepo) GetNotificationConfigsForIncidentUpdate(ctx context.Context, updateID uint) ([]notifications.ViewNotification, error) {
	return m.incidentCfgs, nil
}
func (m *mockNotificationsRepo) GetNotificationConfigsForMaintenanceUpdate(ctx context.Context, updateID uint) ([]notifications.ViewNotification, error) {
	return nil, nil
}
func (m *mockNotificationsRepo) GetIncidentNotificationDetails(ctx context.Context, updateID uint) (notifications.IncidentNotificationDetails, error) {
	return m.incDetails, nil
}
func (m *mockNotificationsRepo) GetMaintenanceNotificationDetails(ctx context.Context, updateID uint) (notifications.MaintenanceNotificationDetails, error) {
	return notifications.MaintenanceNotificationDetails{}, nil
}

// Mock queue.Queue implementation
type mockQueue struct {
	mu       sync.Mutex
	messages []queue.Message
	deleted  []string
	archived []string
}

func (mq *mockQueue) Send(ctx context.Context, queueName string, payload json.RawMessage) (string, error) {
	return "", nil
}
func (mq *mockQueue) SendBatch(ctx context.Context, queueName string, payloads []json.RawMessage) ([]string, error) {
	return nil, nil
}
func (mq *mockQueue) Read(ctx context.Context, queueName string, vt int, limit int) ([]queue.Message, error) {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	res := mq.messages
	mq.messages = nil // empty it after reading to simulate dequeue
	return res, nil
}
func (mq *mockQueue) Delete(ctx context.Context, queueName string, messageID string) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.deleted = append(mq.deleted, messageID)
	return nil
}
func (mq *mockQueue) Archive(ctx context.Context, queueName string, messageID string) error {
	mq.mu.Lock()
	defer mq.mu.Unlock()
	mq.archived = append(mq.archived, messageID)
	return nil
}

func TestDispatcherProcessBatch(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	t.Run("successful incident update dispatch - first message", func(t *testing.T) {
		repo := &mockNotificationsRepo{
			deliveries: make(map[string]notifications.NotificationDelivery),
			incidentCfgs: []notifications.ViewNotification{
				{ID: 101, ViewID: 1, Type: notifications.NotificationTypeSlack, Config: json.RawMessage(`{}`)},
			},
			incDetails: notifications.IncidentNotificationDetails{
				IncidentID:  55,
				UpdateID:    9,
				Title:       "Outage",
				ServiceName: "API",
			},
		}

		payloadBytes, _ := json.Marshal(queue.AlertPayload{
			EventType: queue.EventTypeIncidentUpdate,
			EventID:   9,
		})

		q := &mockQueue{
			messages: []queue.Message{
				{ID: "msg-123", Payload: payloadBytes},
			},
		}

		dispatcher := NewDispatcherApplication(q, repo, logger)
		count, err := dispatcher.processBatch(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		assert.Contains(t, q.deleted, "msg-123")
		assert.Empty(t, q.archived)

		// Verify that a delivery state is stored as "first alert"
		deliveryKey := fmt.Sprintf("%d-%s-%d", 101, "incident", 55)
		repo.mu.Lock()
		d, exists := repo.deliveries[deliveryKey]
		repo.mu.Unlock()
		assert.True(t, exists)
		assert.Equal(t, uint(9), d.LastUpdateID)
		assert.Equal(t, "mock-external-id", d.ExternalIdentifier)
	})

	t.Run("subsequent incident update dispatch - update delivery state", func(t *testing.T) {
		repo := &mockNotificationsRepo{
			deliveries: map[string]notifications.NotificationDelivery{
				fmt.Sprintf("%d-%s-%d", 101, "incident", 55): {
					ID:                 999,
					ViewNotificationID: 101,
					AlertType:          "incident",
					AlertID:            55,
					LastUpdateID:       9,
					ExternalIdentifier: "mock-external-id",
				},
			},
			incidentCfgs: []notifications.ViewNotification{
				{ID: 101, ViewID: 1, Type: notifications.NotificationTypeSlack, Config: json.RawMessage(`{}`)},
			},
			incDetails: notifications.IncidentNotificationDetails{
				IncidentID:  55,
				UpdateID:    10,
				Title:       "Outage",
				ServiceName: "API",
			},
		}

		payloadBytes, _ := json.Marshal(queue.AlertPayload{
			EventType: queue.EventTypeIncidentUpdate,
			EventID:   10,
		})

		q := &mockQueue{
			messages: []queue.Message{
				{ID: "msg-124", Payload: payloadBytes},
			},
		}

		dispatcher := NewDispatcherApplication(q, repo, logger)
		count, err := dispatcher.processBatch(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		assert.Contains(t, q.deleted, "msg-124")

		// Verify the delivery state is updated to last update id 10
		deliveryKey := fmt.Sprintf("%d-%s-%d", 101, "incident", 55)
		repo.mu.Lock()
		d, exists := repo.deliveries[deliveryKey]
		repo.mu.Unlock()
		assert.True(t, exists)
		assert.Equal(t, uint(10), d.LastUpdateID)
	})

	t.Run("stale update skipped - idempotency guard", func(t *testing.T) {
		repo := &mockNotificationsRepo{
			deliveries: map[string]notifications.NotificationDelivery{
				fmt.Sprintf("%d-%s-%d", 101, "incident", 55): {
					ID:                 999,
					ViewNotificationID: 101,
					AlertType:          "incident",
					AlertID:            55,
					LastUpdateID:       10,
					ExternalIdentifier: "mock-external-id",
				},
			},
			incidentCfgs: []notifications.ViewNotification{
				{ID: 101, ViewID: 1, Type: notifications.NotificationTypeSlack, Config: json.RawMessage(`{}`)},
			},
			incDetails: notifications.IncidentNotificationDetails{
				IncidentID:  55,
				UpdateID:    9, // Stale update
				Title:       "Outage",
				ServiceName: "API",
			},
		}

		payloadBytes, _ := json.Marshal(queue.AlertPayload{
			EventType: queue.EventTypeIncidentUpdate,
			EventID:   9,
		})

		q := &mockQueue{
			messages: []queue.Message{
				{ID: "msg-125", Payload: payloadBytes},
			},
		}

		dispatcher := NewDispatcherApplication(q, repo, logger)
		count, err := dispatcher.processBatch(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		assert.Contains(t, q.deleted, "msg-125")

		// Verify delivery state remains at last update id 10
		deliveryKey := fmt.Sprintf("%d-%s-%d", 101, "incident", 55)
		repo.mu.Lock()
		d, exists := repo.deliveries[deliveryKey]
		repo.mu.Unlock()
		assert.True(t, exists)
		assert.Equal(t, uint(10), d.LastUpdateID)
	})

	t.Run("corrupt message archived", func(t *testing.T) {
		repo := &mockNotificationsRepo{}
		q := &mockQueue{
			messages: []queue.Message{
				{ID: "msg-corrupt", Payload: json.RawMessage(`{invalid_json}`)},
			},
		}

		dispatcher := NewDispatcherApplication(q, repo, logger)
		count, err := dispatcher.processBatch(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		assert.Contains(t, q.archived, "msg-corrupt")
		assert.Empty(t, q.deleted)
	})
}
