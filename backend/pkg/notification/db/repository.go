package db

import (
	"context"

	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/mongo"

	notification "DoAn"
	"DoAn/entity"
)

type repo struct {
	collection *mongo.Collection
	logger     log.Logger
}

func (r repo) GetNotification(ctx context.Context, userId int) ([]entity.Notification, error) {
	var resp []entity.Notification
	cursor, err := r.collection.Find(ctx, map[string]interface{}{"user_id": userId})
	if err != nil {
		r.logger.Log("error", err)
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var noti entity.Notification
		err := cursor.Decode(&noti)
		if err != nil {
			r.logger.Log("error", err)
			return nil, err
		}
		resp = append(resp, noti)
	}
	return resp, nil
}

func NewRepository(collection *mongo.Collection, logger log.Logger) notification.NotificationRepository {
	return &repo{
		collection: collection,
		logger:     logger,
	}
}
