package entity

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId    int                `json:"user_id" bson:"user_id"`
	Title     string             `json:"title" bson:"title"`
	Message   string             `json:"message" bson:"message"`
	Type      string             `json:"type" bson:"type"`
	Timestamp int64              `json:"timestamp" bson:"timestamp"`
	RawData   json.RawMessage    `json:"raw_data" bson:"raw_data"`
	IsRead    bool               `json:"is_read" bson:"is_read"`
}
