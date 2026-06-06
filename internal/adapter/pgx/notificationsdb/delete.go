package notificationsdb

import (
	"context"
	_ "embed"
	"fmt"
)

//go:embed queries/delete_view_notification.sql
var deleteViewNotificationQuery string

// Delete view notification destination soft/hard delete
func (r *PostgresNotificationsRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.writeDB.Exec(ctx, deleteViewNotificationQuery, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete view notification: %w", err)
	}
	return nil
}
