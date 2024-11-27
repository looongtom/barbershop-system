package service

import (
	"DoAn/api"
	"context"

	"github.com/go-kit/kit/log"

	notification "DoAn"
)

type NotificationStruct struct {
	repository notification.NotificationRepository
	logger     log.Logger
}

func (n NotificationStruct) GetNotification(ctx context.Context, request api.GetListNotificationRequest) (interface{}, error) {
	listNotification, err := n.repository.GetNotification(ctx, request)
	if err != nil {
		return nil, err
	}
	return listNotification, nil
}

func NewService(repo notification.NotificationRepository, logger log.Logger) notification.NotificationService {
	return &NotificationStruct{
		repository: repo,
		logger:     logger,
	}
}
