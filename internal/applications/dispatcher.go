package applications

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/yash492/statusy/internal/common/apperrors"
	"github.com/yash492/statusy/internal/common/queue"
	"github.com/yash492/statusy/internal/domain/notifications"
	"golang.org/x/sync/errgroup"
)

type dispatchFunc func(ctx context.Context, eventID uint) error

type DispatcherApplication struct {
	queue        queue.Queue
	viewsRepo    notifications.NotificationsRepository
	notifier     notifications.Notifier
	lg           *slog.Logger
	pollInterval time.Duration
	handlers     map[queue.EventType]dispatchFunc
}

func NewDispatcherApplication(
	q queue.Queue,
	repo notifications.NotificationsRepository,
	notifier notifications.Notifier,
	lg *slog.Logger,
) *DispatcherApplication {
	d := &DispatcherApplication{
		queue:        q,
		viewsRepo:    repo,
		notifier:     notifier,
		lg:           lg,
		pollInterval: 2 * time.Second,
	}
	d.handlers = map[queue.EventType]dispatchFunc{
		queue.EventTypeIncidentUpdate:    d.dispatchIncidentUpdate,
		queue.EventTypeMaintenanceUpdate: d.dispatchMaintenanceUpdate,
	}
	return d
}

// Start runs the background worker loop polling the queue with conditional backoff and smart shutdown
func (d *DispatcherApplication) Start(ctx context.Context) error {
	d.lg.InfoContext(ctx, "starting notification dispatcher background worker")
	for {
		select {
		case <-ctx.Done():
			d.lg.InfoContext(ctx, "notification dispatcher stopped gracefully")
			return ctx.Err()
		default:
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

// processBatch reads and handles messages sequentially to guarantee strict execution ordering
func (d *DispatcherApplication) processBatch(ctx context.Context) (int, error) {
	messages, err := d.queue.Read(ctx, "notifications", 30, 10)
	if err != nil {
		return 0, fmt.Errorf("failed to read from queue: %w", err)
	}

	if len(messages) == 0 {
		return 0, nil
	}

	d.lg.DebugContext(ctx, "processing queue batch", slog.Int("batch_size", len(messages)))

	for _, msg := range messages {
		envelope, err := queue.UnmarshalMessage[queue.AlertPayload](msg)
		if err != nil {
			d.lg.ErrorContext(ctx, "corrupt message payload, archiving", slog.String("msg_id", msg.ID), slog.Any("err", err))
			_ = d.queue.Archive(ctx, "notifications", msg.ID)
			continue
		}

		handler, ok := d.handlers[envelope.Payload.EventType]
		if !ok {
			d.lg.WarnContext(ctx, "unknown event type, archiving", slog.String("type", string(envelope.Payload.EventType)))
			_ = d.queue.Archive(ctx, "notifications", msg.ID)
			continue
		}

		// Strict sequential execution: wait for all channels of this update to complete
		if err := handler(ctx, envelope.Payload.EventID); err != nil {
			return len(messages), fmt.Errorf("dispatch failed for msg %s: %w", msg.ID, err)
		}

		if err := d.queue.Delete(ctx, "notifications", msg.ID); err != nil {
			return len(messages), fmt.Errorf("failed to delete msg %s: %w", msg.ID, err)
		}
	}

	return len(messages), nil
}

func (d *DispatcherApplication) dispatchIncidentUpdate(ctx context.Context, updateID uint) error {
	channels, err := d.viewsRepo.GetNotificationConfigsForIncidentUpdate(ctx, updateID)
	if err != nil {
		return fmt.Errorf("failed to resolve incident notification configs: %w", err)
	}

	details, err := d.viewsRepo.GetIncidentNotificationDetails(ctx, updateID)
	if err != nil {
		return fmt.Errorf("failed to get incident notification details: %w", err)
	}

	errGroup, errGroupCtx := errgroup.WithContext(ctx)
	errGroup.SetLimit(10)
	for _, channel := range channels {
		ch := channel
		errGroup.Go(func() error {
			delivery, err := d.viewsRepo.GetDelivery(errGroupCtx, ch.ID, notifications.AlertTypeIncident, details.IncidentID)

			isFirst := false
			if err != nil {
				if err, ok := errors.AsType[*apperrors.AppError](err); ok {
					if err.Type == apperrors.TypeNotFound {
						isFirst = true
					} else {
						return fmt.Errorf("failed to get delivery state: %w", err)
					}
				}
			}

			if !isFirst && delivery.LastUpdateID >= updateID {
				d.lg.DebugContext(errGroupCtx, "incident notification already delivered", slog.Uint64("channel_id", uint64(ch.ID)), slog.Uint64("update_id", uint64(updateID)))
				return nil
			}

			extID, sendErr := d.notifier.SendIncident(errGroupCtx, ch, isFirst, details, delivery.ExternalIdentifier)
			if sendErr != nil {
				d.lg.ErrorContext(errGroupCtx, "failed to dispatch notification",
					slog.Uint64("channel_id", uint64(ch.ID)),
					slog.String("channel_type", string(ch.Type)),
					slog.Any("err", sendErr),
				)
				dbErr := d.viewsRepo.SaveDeliveryFailure(errGroupCtx, notifications.NotificationDeliveryFailure{
					ViewNotificationID: ch.ID,
					AlertType:          notifications.AlertTypeIncident,
					AlertID:            details.IncidentID,
					UpdateID:           updateID,
					ErrorMessage:       sendErr.Error(),
				})
				if dbErr != nil {
					d.lg.ErrorContext(errGroupCtx, "failed to save delivery failure to db", slog.Any("err", dbErr))
				}
				return nil
			}

			if isFirst {
				err = d.viewsRepo.SaveDelivery(errGroupCtx, notifications.NotificationDelivery{
					ViewNotificationID: ch.ID,
					AlertType:          notifications.AlertTypeIncident,
					AlertID:            details.IncidentID,
					LastUpdateID:       updateID,
					ExternalIdentifier: extID,
				})
			} else {
				err = d.viewsRepo.UpdateDelivery(errGroupCtx, delivery.ID, updateID, extID)
			}
			if err != nil {
				return fmt.Errorf("failed to record delivery: %w", err)
			}
			return nil
		})
	}
	return errGroup.Wait()
}

func (d *DispatcherApplication) dispatchMaintenanceUpdate(ctx context.Context, updateID uint) error {
	channels, err := d.viewsRepo.GetNotificationConfigsForMaintenanceUpdate(ctx, updateID)
	if err != nil {
		return fmt.Errorf("failed to resolve maintenance notification configs: %w", err)
	}

	details, err := d.viewsRepo.GetMaintenanceNotificationDetails(ctx, updateID)
	if err != nil {
		return fmt.Errorf("failed to get maintenance notification details: %w", err)
	}

	g, gCtx := errgroup.WithContext(ctx)
	for _, channel := range channels {
		ch := channel
		g.Go(func() error {
			delivery, err := d.viewsRepo.GetDelivery(gCtx, ch.ID, notifications.AlertTypeScheduledMaintenance, details.MaintenanceID)

			isFirst := false
			if err != nil {
				var appErr *apperrors.AppError
				if errors.As(err, &appErr) && appErr.Type == apperrors.TypeNotFound {
					isFirst = true
				} else {
					return fmt.Errorf("failed to get delivery state: %w", err)
				}
			}

			if !isFirst && delivery.LastUpdateID >= updateID {
				d.lg.DebugContext(gCtx, "maintenance notification already delivered", slog.Uint64("channel_id", uint64(ch.ID)), slog.Uint64("update_id", uint64(updateID)))
				return nil
			}

			extID, sendErr := d.notifier.SendMaintenance(gCtx, ch, isFirst, details, delivery.ExternalIdentifier)
			if sendErr != nil {
				d.lg.ErrorContext(gCtx, "failed to dispatch notification",
					slog.Uint64("channel_id", uint64(ch.ID)),
					slog.String("channel_type", string(ch.Type)),
					slog.Any("err", sendErr),
				)
				dbErr := d.viewsRepo.SaveDeliveryFailure(gCtx, notifications.NotificationDeliveryFailure{
					ViewNotificationID: ch.ID,
					AlertType:          notifications.AlertTypeScheduledMaintenance,
					AlertID:            details.MaintenanceID,
					UpdateID:           updateID,
					ErrorMessage:       sendErr.Error(),
				})
				if dbErr != nil {
					d.lg.ErrorContext(gCtx, "failed to save delivery failure to db", slog.Any("err", dbErr))
				}
				return nil
			}

			if isFirst {
				err = d.viewsRepo.SaveDelivery(gCtx, notifications.NotificationDelivery{
					ViewNotificationID: ch.ID,
					AlertType:          notifications.AlertTypeScheduledMaintenance,
					AlertID:            details.MaintenanceID,
					LastUpdateID:       updateID,
					ExternalIdentifier: extID,
				})
			} else {
				err = d.viewsRepo.UpdateDelivery(gCtx, delivery.ID, updateID, extID)
			}
			if err != nil {
				return fmt.Errorf("failed to record delivery: %w", err)
			}
			return nil
		})
	}
	return g.Wait()
}
