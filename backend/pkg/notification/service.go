package notification

import (
	"context"
)

type NotificationService interface {
	GetNotification(ctx context.Context, userId int) (interface{}, error)
}
