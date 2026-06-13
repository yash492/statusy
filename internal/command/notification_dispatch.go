package command

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/common/queue"
	"github.com/yash492/statusy/internal/domain/notifications"
	"golang.org/x/sync/errgroup"
)

const notificationRetryCountLimit = 3

type NotificationDispatchCmd struct {
	queue     queue.Queue
	viewsRepo notifications.NotificationsRepository
	notifier  notifications.Notifier
	lg        *slog.Logger
	handlers  map[queue.EventType]func(ctx context.Context, eventID uint) error
}

func NewNotificationDispatchCmd(
	q queue.Queue,
	viewsRepo notifications.NotificationsRepository,
	notifier notifications.Notifier,
	lg *slog.Logger,
) NotificationDispatchCmd {
	cmd := NotificationDispatchCmd{
		queue:     q,
		viewsRepo: viewsRepo,
		notifier:  notifier,
		lg:        lg,
	}
	cmd.handlers = map[queue.EventType]func(ctx context.Context, eventID uint) error{
		queue.EventTypeIncidentUpdate:    cmd.dispatchIncidentUpdate,
		queue.EventTypeMaintenanceUpdate: cmd.dispatchMaintenanceUpdate,
	}
	return cmd
}

// Execute reads and handles messages sequentially to guarantee strict execution ordering.
// Returns the number of messages read and any processing error.
func (c NotificationDispatchCmd) Execute(ctx context.Context) (int, error) {
	messages, err := c.queue.Read(ctx, queue.Notifications, 30, 10)
	if err != nil {
		return 0, fmt.Errorf("failed to read from queue: %w", err)
	}

	if len(messages) == 0 {
		return 0, nil
	}

	c.lg.DebugContext(ctx, "processing queue batch", slog.Int("batch_size", len(messages)))

	for _, msg := range messages {
		envelope, err := queue.UnmarshalMessage[queue.AlertPayload](msg)
		if err != nil {
			c.lg.ErrorContext(ctx, "corrupt message payload, archiving", slog.String("msg_id", msg.ID), slog.Any("err", err))
			_ = c.queue.Archive(ctx, queue.Notifications, msg.ID)
			continue
		}

		handler, ok := c.handlers[envelope.Payload.EventType]
		if !ok {
			c.lg.WarnContext(ctx, "unknown event type, archiving", slog.String("type", string(envelope.Payload.EventType)))
			_ = c.queue.Archive(ctx, queue.Notifications, msg.ID)
			continue
		}

		if err := handler(ctx, envelope.Payload.EventID); err != nil {
			if msg.ReadCount <= notificationRetryCountLimit {
				return len(messages), fmt.Errorf("dispatch failed for msg %s: %w", msg.ID, err)
			}
			c.lg.ErrorContext(ctx, "dispatch failed after max retries, archiving", slog.String("msg_id", msg.ID), slog.Int("read_count", msg.ReadCount), slog.Any("err", err))
			_ = c.queue.Archive(ctx, queue.Notifications, msg.ID)
			continue
		}

		if err := c.queue.Delete(ctx, queue.Notifications, msg.ID); err != nil {
			return len(messages), fmt.Errorf("failed to delete msg %s: %w", msg.ID, err)
		}
	}

	return len(messages), nil
}

func (c NotificationDispatchCmd) dispatchIncidentUpdate(ctx context.Context, updateID uint) error {
	channels, err := c.viewsRepo.GetNotificationConfigsForIncidentUpdate(ctx, updateID)
	if err != nil {
		return fmt.Errorf("failed to resolve incident notification configs: %w", err)
	}

	details, err := c.viewsRepo.GetIncidentNotificationDetails(ctx, updateID)
	if err != nil {
		return fmt.Errorf("failed to get incident notification details: %w", err)
	}

	errGroup, errGroupCtx := errgroup.WithContext(ctx)
	errGroup.SetLimit(10)
	for _, channel := range channels {
		ch := channel
		errGroup.Go(func() error {
			delivery, err := c.viewsRepo.GetDelivery(errGroupCtx, ch.ID, notifications.AlertTypeIncident, details.IncidentID)

			notFound := false
			if err, ok := errors.AsType[*apperrors.AppError](err); ok {
				if err.Type == apperrors.TypeNotFound {
					notFound = true
				} else {
					return fmt.Errorf("failed to get delivery state: %w", err)
				}
			}

			if !notFound && delivery.LastUpdateID >= updateID {
				c.lg.DebugContext(errGroupCtx, "incident notification already delivered", slog.Uint64("channel_id", uint64(ch.ID)), slog.Uint64("update_id", uint64(updateID)))
				return nil
			}

			extID, sendErr := c.notifier.SendIncident(errGroupCtx, ch, details, delivery.ExternalIdentifier)
			if sendErr != nil {
				c.lg.ErrorContext(errGroupCtx, "failed to dispatch notification",
					slog.Uint64("channel_id", uint64(ch.ID)),
					slog.String("channel_type", string(ch.Type)),
					slog.Any("err", sendErr),
				)
				dbErr := c.viewsRepo.SaveDeliveryFailure(errGroupCtx, notifications.NotificationDeliveryFailure{
					ViewNotificationID: ch.ID,
					AlertType:          notifications.AlertTypeIncident,
					AlertID:            details.IncidentID,
					UpdateID:           updateID,
					ErrorMessage:       sendErr.Error(),
				})
				if dbErr != nil {
					c.lg.ErrorContext(errGroupCtx, "failed to save delivery failure to db", slog.Any("err", dbErr))
				}
				return nil
			}

			err = c.viewsRepo.UpsertDelivery(errGroupCtx, notifications.NotificationDelivery{
				ViewNotificationID: ch.ID,
				AlertType:          notifications.AlertTypeIncident,
				AlertID:            details.IncidentID,
				LastUpdateID:       updateID,
				ExternalIdentifier: extID,
			})
			if err != nil {
				return fmt.Errorf("failed to record delivery: %w", err)
			}
			return nil
		})
	}
	return errGroup.Wait()
}

func (c NotificationDispatchCmd) dispatchMaintenanceUpdate(ctx context.Context, updateID uint) error {
	channels, err := c.viewsRepo.GetNotificationConfigsForMaintenanceUpdate(ctx, updateID)
	if err != nil {
		return fmt.Errorf("failed to resolve maintenance notification configs: %w", err)
	}

	details, err := c.viewsRepo.GetMaintenanceNotificationDetails(ctx, updateID)
	if err != nil {
		return fmt.Errorf("failed to get maintenance notification details: %w", err)
	}

	g, gCtx := errgroup.WithContext(ctx)
	for _, channel := range channels {
		ch := channel
		g.Go(func() error {
			delivery, err := c.viewsRepo.GetDelivery(gCtx, ch.ID, notifications.AlertTypeScheduledMaintenance, details.MaintenanceID)

			notFound := false
			if err, ok := errors.AsType[*apperrors.AppError](err); ok {
				if err.Type == apperrors.TypeNotFound {
					notFound = true
				} else {
					return fmt.Errorf("failed to get delivery state: %w", err)
				}
			}

			if !notFound && delivery.LastUpdateID >= updateID {
				c.lg.DebugContext(gCtx, "maintenance notification already delivered", slog.Uint64("channel_id", uint64(ch.ID)), slog.Uint64("update_id", uint64(updateID)))
				return nil
			}

			extID, sendErr := c.notifier.SendMaintenance(gCtx, ch, details, delivery.ExternalIdentifier)
			if sendErr != nil {
				c.lg.ErrorContext(gCtx, "failed to dispatch notification",
					slog.Uint64("channel_id", uint64(ch.ID)),
					slog.String("channel_type", string(ch.Type)),
					slog.Any("err", sendErr),
				)
				dbErr := c.viewsRepo.SaveDeliveryFailure(gCtx, notifications.NotificationDeliveryFailure{
					ViewNotificationID: ch.ID,
					AlertType:          notifications.AlertTypeScheduledMaintenance,
					AlertID:            details.MaintenanceID,
					UpdateID:           updateID,
					ErrorMessage:       sendErr.Error(),
				})
				if dbErr != nil {
					c.lg.ErrorContext(gCtx, "failed to save delivery failure to db", slog.Any("err", dbErr))
				}
				return nil
			}

			err = c.viewsRepo.UpsertDelivery(gCtx, notifications.NotificationDelivery{
				ViewNotificationID: ch.ID,
				AlertType:          notifications.AlertTypeScheduledMaintenance,
				AlertID:            details.MaintenanceID,
				LastUpdateID:       updateID,
				ExternalIdentifier: extID,
			})
			if err != nil {
				return fmt.Errorf("failed to record delivery: %w", err)
			}
			return nil
		})
	}
	return g.Wait()
}
