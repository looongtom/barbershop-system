package notification

import (
	"DoAn/api"
	"context"
)

type NotificationService interface {
	GetNotification(ctx context.Context, request api.GetListNotificationRequest) (interface{}, error)
}
