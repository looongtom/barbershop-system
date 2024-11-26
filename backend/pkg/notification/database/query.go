package database

import (
	"context"
	"os"

	"DoAn/entity"
)

func SaveNotification(notification entity.Notification) error {
	// save the notification to the database
	collectionMongo := ConnectMongo(os.Getenv("NotificationCollectionMongo"))
	_, err := collectionMongo.InsertOne(context.Background(), notification)
	if err != nil {
		return err
	}
	return nil
}
