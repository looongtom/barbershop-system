package notification

import (
	"context"

	"DoAn/entity"
)

type NotificationRepository interface {
	GetNotification(ctx context.Context, userId int) ([]entity.Notification, error)
}
