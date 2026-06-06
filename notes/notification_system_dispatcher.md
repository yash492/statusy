# PGMQ-Based Notification System and Dispatcher

Design and implement a notification subsystem using PostgreSQL Message Queue (`pgmq`). The notification dispatcher will run as a background worker in a separate Go routine managed by the main application's lifecycle context.

## Design Decisions Summary (User Approved)

1. **Route on Dequeue (Strategy B)**:
   - When a new update is saved, we push a single message referencing that update to the queue: `{"event_type": "incident_update", "update_id": 123}` or `{"event_type": "maintenance_update", "update_id": 456}`.
   - The background dispatcher worker reads this message, queries the database to resolve the target views/integrations, and posts to each destination channel.
2. **Polymorphic Delivery Mapping**:
   - To update/edit existing chat messages (Slack, Discord) and resolve incidents (PagerDuty, SolarwindsIncidentResponse), we map local alerts to external message keys.
   - We will use a polymorphic `notification_deliveries` table with `alert_type` ('incident' or 'sm'), `alert_id`, and `last_update_id` columns.
3. **Queue Payload**: Enqueue specific update IDs: `incident_updates.id` and `scheduled_maintenance_updates.id`. This inherently covers both brand new incidents (first update) and subsequent updates.
4. **Normalized Database Storage**: Store integration configurations in a `view_notifications` table linked to `views`.
5. **No External Library**: Call standard PGMQ SQL functions via `pgx` directly.
6. **Docker Database Image**: Update `compose.yaml` to `ghcr.io/pgmq/pg18-pgmq:latest`.
7. **Concurrent Dispatch**: Trigger external alerts concurrently in parallel goroutines using `errgroup.Group` to ensure low dispatch latency, while preserving target-level idempotency checks.
8. **Smart Backoff Loop (No Busy-Wait)**: Run a forever-loop worker that instantly loops to process the next message batch if the queue contains items, but backs off and sleeps only when the queue is empty or errors occur.
9. **Bulk Processing**: Read and process queue messages in batches (up to 10 at once) using nested concurrency to process multiple events in parallel, each containing parallel channel dispatches.
10. **Batch Enqueuing**: Perform batch inserts into pgmq from the scraper using `pgmq.send_batch('notifications', ARRAY[payload1, payload2, ...])` to ensure maximum DB performance.

---

## Alert Resolution Design (Strategy B)

### 1. Enqueue Phase (Batch Inserts)
When the scraper finishes saving a batch of updates, it collects the returned IDs of the newly inserted rows (obtained via `ON CONFLICT DO NOTHING RETURNING *`). It then enqueues all of them to PGMQ in a single database batch insert query:

```sql
SELECT * FROM pgmq.send_batch('notifications', $1::jsonb[]);
```

### 2. Dispatch Resolution Queries
When dequeuing an update, the worker queries the database to find target integrations.

#### For Incident Updates:
```sql
SELECT DISTINCT
    vn.id AS notification_id,
    vn.type AS notification_type,
    vn.config AS notification_config,
    vn.view_id AS view_id
FROM view_notifications vn
JOIN views v ON v.id = vn.view_id AND v.deleted_at IS NULL
JOIN view_services vs ON vs.view_id = v.id AND vs.deleted_at IS NULL
JOIN incident_updates iu ON iu.id = :update_id
JOIN incidents i ON i.id = iu.incident_id AND i.service_id = vs.service_id
WHERE vs.monitor_incidents = true
  AND vn.deleted_at IS NULL
  AND (
    vs.include_all_components = true
    
    -- Match if the incident affects no components
    OR NOT EXISTS (SELECT 1 FROM incident_components ic WHERE ic.incident_id = i.id)

    -- Match if the incident affects a component selected in this view
    OR EXISTS (
      SELECT 1 FROM incident_components ic
      JOIN view_service_components vsc ON vsc.component_id = ic.component_id AND vsc.deleted_at IS NULL
      WHERE ic.incident_id = i.id AND vsc.view_service_id = vs.id
    )
    
    -- Match if the incident affects a component in a group selected in this view
    OR EXISTS (
      SELECT 1 FROM incident_components ic
      JOIN components c ON c.id = ic.component_id AND c.deleted_at IS NULL
      JOIN view_service_component_groups vscg ON vscg.component_group_id = c.component_group_id AND vscg.deleted_at IS NULL
      WHERE ic.incident_id = i.id AND vscg.view_service_id = vs.id
    )
  );
```

#### For Scheduled Maintenance Updates:
```sql
SELECT DISTINCT
    vn.id AS notification_id,
    vn.type AS notification_type,
    vn.config AS notification_config,
    vn.view_id AS view_id
FROM view_notifications vn
JOIN views v ON v.id = vn.view_id AND v.deleted_at IS NULL
JOIN view_services vs ON vs.view_id = v.id AND vs.deleted_at IS NULL
JOIN scheduled_maintenance_updates smu ON smu.id = :update_id
JOIN scheduled_maintenances sm ON sm.id = smu.scheduled_maintenance_id AND sm.service_id = vs.service_id
WHERE vs.monitor_scheduled_maintenances = true
  AND vn.deleted_at IS NULL
  AND (
    vs.include_all_components = true
    
    -- Match if the maintenance affects no components
    OR NOT EXISTS (SELECT 1 FROM scheduled_maintenance_components smc WHERE smc.scheduled_maintenance_id = sm.id)

    -- Match if the maintenance affects a component selected in this view
    OR EXISTS (
      SELECT 1 FROM scheduled_maintenance_components smc
      JOIN view_service_components vsc ON vsc.component_id = smc.component_id AND vsc.deleted_at IS NULL
      WHERE smc.scheduled_maintenance_id = sm.id AND vsc.view_service_id = vs.id
    )
    
    -- Match if the maintenance affects a component in a group selected in this view
    OR EXISTS (
      SELECT 1 FROM scheduled_maintenance_components smc
      JOIN components c ON c.id = smc.component_id AND c.deleted_at IS NULL
      JOIN view_service_component_groups vscg ON vscg.component_group_id = c.component_group_id AND vscg.deleted_at IS NULL
      WHERE smc.scheduled_maintenance_id = sm.id AND vscg.view_service_id = vs.id
    )
  );
```

---

## Polymorphic Delivery and Idempotency Schema Design

The `notification_deliveries` table tracks the external message/alert ID dynamically based on alert type.

```sql
CREATE TABLE IF NOT EXISTS notification_deliveries (
    id SERIAL PRIMARY KEY,
    view_notification_id INT NOT NULL REFERENCES view_notifications(id) ON DELETE CASCADE,
    
    -- Polymorphic type mapping
    alert_type TEXT NOT NULL, -- 'incident' or 'sm'
    alert_id INT NOT NULL,    -- Maps to incidents.id or scheduled_maintenances.id
    
    -- Tracks the last update successfully dispatched (Idempotency)
    last_update_id INT NOT NULL, -- Maps to incident_updates.id or scheduled_maintenance_updates.id
    
    -- Delivery key from external API (e.g. Slack ts, PagerDuty dedup_key, Discord message ID)
    external_identifier TEXT NOT NULL,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    
    CONSTRAINT unique_notification_delivery UNIQUE (view_notification_id, alert_type, alert_id)
);
```

---

## Proposed Code Documentation

### 1. Status Normalizer (`internal/common/status_normalizer.go`)
Normalizes raw provider statuses into standard, consistent values.
```go
package common

import "strings"

// NormalizeIncidentStatus standardizes incident states to: investigating, identified, monitoring, resolved
func NormalizeIncidentStatus(status string) string {
	s := strings.ToLower(strings.TrimSpace(status))
	switch s {
	case "investigating":
		return "investigating"
	case "identified", "identified_incident":
		return "identified"
	case "monitoring", "monitoring_incident":
		return "monitoring"
	case "resolved", "completed":
		return "resolved"
	default:
		return "investigating"
	}
}

// NormalizeMaintenanceStatus standardizes scheduled maintenance states to: scheduled, in_progress, verifying, completed
func NormalizeMaintenanceStatus(status string) string {
	s := strings.ToLower(strings.TrimSpace(status))
	switch s {
	case "scheduled", "planned":
		return "scheduled"
	case "in_progress", "active":
		return "in_progress"
	case "verifying", "monitoring":
		return "verifying"
	case "completed", "resolved":
		return "completed"
	default:
		return "scheduled"
	}
}
```

### 2. View Notifications Domain (`internal/domain/views/notification.go`)
```go
package views

import (
	"context"
	"encoding/json"
	"time"
)

type NotificationType string

const (
	NotificationTypeSlack     NotificationType = "slack"
	NotificationTypeDiscord   NotificationType = "discord"
	NotificationTypeMsTeams   NotificationType = "msteams"
	NotificationTypePagerDuty NotificationType = "pagerduty"
	NotificationTypeSolarwindsIncidentResponse NotificationType = "solarwinds_incident_response"
	NotificationTypeWebhook   NotificationType = "webhook"
)

type ViewNotification struct {
	ID        uint             `json:"id"`
	ViewID    uint             `json:"view_id"`
	Type      NotificationType `json:"type"`
	Config    json.RawMessage  `json:"config"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type NotificationDelivery struct {
	ID                 uint      `json:"id"`
	ViewNotificationID uint      `json:"view_notification_id"`
	AlertType          string    `json:"alert_type"` // "incident" or "sm"
	AlertID            uint      `json:"alert_id"`
	LastUpdateID       uint      `json:"last_update_id"`
	ExternalIdentifier string    `json:"external_identifier"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type NotificationsRepository interface {
	Save(ctx context.Context, vn ViewNotification) (ViewNotification, error)
	GetByViewID(ctx context.Context, viewID uint) ([]ViewNotification, error)
	Delete(ctx context.Context, id uint) error
	
	GetDelivery(ctx context.Context, channelID uint, alertType string, alertID uint) (NotificationDelivery, error)
	SaveDelivery(ctx context.Context, delivery NotificationDelivery) error
	UpdateDelivery(ctx context.Context, deliveryID uint, lastUpdateID uint) error
}
```

### 3. Background Dispatcher Worker (`internal/applications/dispatcher.go`)
```go
package applications

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yash492/statusy/internal/domain/views"
	"golang.org/x/sync/errgroup"
)

type QueueMessage struct {
	MsgID     int64           `json:"msg_id"`
	Payload   json.RawMessage `json:"message"`
}

type AlertPayload struct {
	EventType string `json:"event_type"` // "incident_update" or "maintenance_update"
	UpdateID  uint   `json:"update_id"`
}

type DispatcherApplication struct {
	db           *pgxpool.Pool
	viewsRepo    views.NotificationsRepository
	httpClient   *http.Client
	lg           *slog.Logger
	pollInterval time.Duration // Backoff time when queue is empty
}

func NewDispatcherApplication(db *pgxpool.Pool, repo views.NotificationsRepository, lg *slog.Logger) *DispatcherApplication {
	return &DispatcherApplication{
		db:           db,
		viewsRepo:    repo,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		lg:           lg,
		pollInterval: 2 * time.Second,
	}
}

// Start runs the forever worker loop with smart backoff
func (d *DispatcherApplication) Start(ctx context.Context) error {
	d.lg.Info("starting notification dispatcher background worker")

	for {
		select {
		case <-ctx.Done():
			d.lg.Info("notification dispatcher stopped gracefully")
			return ctx.Err()
		default:
			// Attempt to process a batch of messages from PGMQ
			processedCount, err := d.processBatch(ctx)
			if err != nil {
				d.lg.Error("error processing queue batch", slog.Any("err", err))
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(d.pollInterval):
				}
				continue
			}

			if processedCount == 0 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(d.pollInterval):
				}
			}
		}
	}
}

// processBatch reads up to 10 messages from PGMQ and dispatches them concurrently
func (d *DispatcherApplication) processBatch(ctx context.Context) (int, error) {
	rows, err := d.db.Query(ctx, "SELECT msg_id, message FROM pgmq.read('notifications', 30, 10);")
	if err != nil {
		return 0, fmt.Errorf("failed to read from pgmq queue: %w", err)
	}
	defer rows.Close()

	var messages []QueueMessage
	for rows.Next() {
		var msg QueueMessage
		var rawPayload []byte
		if err := rows.Scan(&msg.MsgID, &rawPayload); err != nil {
			return 0, fmt.Errorf("failed to scan batch queue row: %w", err)
		}
		msg.Payload = rawPayload
		messages = append(messages, msg)
	}

	if len(messages) == 0 {
		return 0, nil
	}

	d.lg.DebugContext(ctx, "processing pgmq batch", slog.Int("batch_size", len(messages)))

	g, gCtx := errgroup.WithContext(ctx)

	for _, msg := range messages {
		m := msg
		g.Go(func() error {
			var payload AlertPayload
			if err := json.Unmarshal(m.Payload, &payload); err != nil {
				d.lg.ErrorContext(gCtx, "corrupt message payload, archiving", slog.Int64("msg_id", m.MsgID), slog.Any("err", err))
				_, _ = d.db.Exec(gCtx, "SELECT pgmq.archive('notifications', $1);", m.MsgID)
				return nil
			}

			var dispatchErr error
			switch payload.EventType {
			case "incident_update":
				dispatchErr = d.dispatchIncidentUpdate(gCtx, payload.UpdateID)
			case "maintenance_update":
				dispatchErr = d.dispatchMaintenanceUpdate(gCtx, payload.UpdateID)
			default:
				d.lg.WarnContext(gCtx, "unknown event type, archiving", slog.String("type", payload.EventType))
				_, _ = d.db.Exec(gCtx, "SELECT pgmq.archive('notifications', $1);", m.MsgID)
				return nil
			}

			if dispatchErr != nil {
				return fmt.Errorf("dispatch failed for msg %d: %w", m.MsgID, dispatchErr)
			}

			_, err = d.db.Exec(gCtx, "SELECT pgmq.delete('notifications', $1);", m.MsgID)
			if err != nil {
				return fmt.Errorf("failed to delete msg %d: %w", m.MsgID, err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return len(messages), err
	}

	return len(messages), nil
}

func (d *DispatcherApplication) dispatchIncidentUpdate(ctx context.Context, updateID uint) error {
	channels, err := d.viewsRepo.GetNotificationConfigsForIncidentUpdate(ctx, updateID)
	if err != nil {
		return fmt.Errorf("failed to resolve view channels: %w", err)
	}

	incidentData, err := d.viewsRepo.GetIncidentNotificationDetails(ctx, updateID)
	if err != nil {
		return fmt.Errorf("failed to get incident notification details: %w", err)
	}

	g, gCtx := errgroup.WithContext(ctx)

	for _, channel := range channels {
		ch := channel
		g.Go(func() error {
			delivery, err := d.viewsRepo.GetDelivery(gCtx, ch.ID, "incident", incidentData.IncidentID)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("failed to fetch delivery state for channel %d: %w", ch.ID, err)
			}

			if err == nil && delivery.LastUpdateID >= updateID {
				return nil
			}

			isFirstMessage := errors.Is(err, pgx.ErrNoRows)
			var extID string
			var dispatchErr error

			switch ch.Type {
			case views.NotificationTypeSlack:
				extID, dispatchErr = d.sendSlackAlert(gCtx, ch.Config, incidentData, delivery.ExternalIdentifier, isFirstMessage)
			case views.NotificationTypeDiscord:
				extID, dispatchErr = d.sendDiscordAlert(gCtx, ch.Config, incidentData, delivery.ExternalIdentifier, isFirstMessage)
			case views.NotificationTypePagerDuty:
				extID, dispatchErr = d.sendPagerDutyAlert(gCtx, ch.Config, incidentData, delivery.ExternalIdentifier, isFirstMessage)
			// ... other integrations ...
			}

			if dispatchErr != nil {
				return fmt.Errorf("dispatch failed for channel %d: %w", ch.ID, dispatchErr)
			}

			if isFirstMessage {
				err = d.viewsRepo.SaveDelivery(gCtx, views.NotificationDelivery{
					ViewNotificationID: ch.ID,
					AlertType:          "incident",
					AlertID:            incidentData.IncidentID,
					LastUpdateID:       updateID,
					ExternalIdentifier: extID,
				})
			} else {
				err = d.viewsRepo.UpdateDelivery(gCtx, delivery.ID, updateID)
			}
			if err != nil {
				return fmt.Errorf("failed to update delivery tracking for channel %d: %w", ch.ID, err)
			}

			return nil
		})
	}

	return g.Wait()
}
```
