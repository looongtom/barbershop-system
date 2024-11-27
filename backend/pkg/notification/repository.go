package notification

import (
	"DoAn/api"
	"context"

	"DoAn/entity"
)

type NotificationRepository interface {
	GetNotification(ctx context.Context, req api.GetListNotificationRequest) ([]entity.Notification, error)
}
