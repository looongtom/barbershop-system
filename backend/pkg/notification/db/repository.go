package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"DoAn/api"

	"github.com/go-kit/kit/log"
	"go.mongodb.org/mongo-driver/mongo"

	notification "DoAn"
	"DoAn/entity"
)

type repo struct {
	collection *mongo.Collection
	logger     log.Logger
}

func (r repo) GetNotification(ctx context.Context, req api.GetListNotificationRequest) ([]entity.Notification, error) {
	var resp []entity.Notification
	cursor, err := r.collection.Find(ctx, map[string]interface{}{
		"user_id": req.UserId,
		"timestamp": map[string]interface{}{
			"$gte": req.StartTime,
			"$lte": req.EndTime,
		},
	}, options.Find().SetSort(map[string]interface{}{
		"timestamp": -1,
	}))
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
