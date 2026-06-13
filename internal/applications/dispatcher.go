package applications

import (
	"context"
	"log/slog"
	"time"

	"github.com/yash492/statusy/internal/command"
	"github.com/yash492/statusy/internal/common/queue"
	"github.com/yash492/statusy/internal/domain/notifications"
)

type DispatcherApplication struct {
	notificationDispatch command.NotificationDispatchCmd
	lg                   *slog.Logger
	pollInterval         time.Duration
}

func NewDispatcherApplication(
	q queue.Queue,
	repo notifications.NotificationsRepository,
	notifier notifications.Notifier,
	lg *slog.Logger,
) *DispatcherApplication {
	return &DispatcherApplication{
		notificationDispatch: command.NewNotificationDispatchCmd(q, repo, notifier, lg),
		lg:                   lg,
		pollInterval:         2 * time.Second,
	}
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
			processedCount, err := d.notificationDispatch.Execute(ctx)
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
