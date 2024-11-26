package api

import (
	"mime/multipart"
)

type GetListNotificationRequest struct {
	UserId    int   `json:"user_id"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

type CreatePreviewImageRequest struct {
	Url       multipart.File `json:"url"`
	AccountId int            `json:"account_id"`
}

type UpdateImageRequest struct {
	SelfImg   multipart.File `json:"selfImg"`
	ShapeImg  multipart.File `json:"shapeImg"`
	ColorImg  multipart.File `json:"colorImg"`
	AccountId int            `json:"account_id"`
}

type KafkaPreviewImageRequest struct {
	ID        int    `json:"id,omitempty"`
	Url       string `json:"url"`
	CreatedAt int64  `json:"created_at,omitempty"`
	AccountId int    `json:"account_id"`
	SelfImg   string `json:"self_img"`
	ShapeImg  string `json:"shape_img"`
	ColorImg  string `json:"color_img"`
}

type HairFastResult struct {
	SelfImgCloud      string `json:"self_img_cloud"`
	ShapeImgCloud     string `json:"shape_img_cloud"`
	ColorImgCloud     string `json:"color_img_cloud"`
	GeneratedImgCloud string `json:"generated_img_cloud"`
}

type KafkaBookingRequest struct {
	UUID          string  `json:"uuid"`
	CustomerID    int     `json:"customer_id"`
	BarberId      int     `json:"barber_id"`
	ResultId      int     `json:"result_id,omitempty"`
	Status        string  `json:"status"`
	Price         float32 `json:"price"`
	SlotId        int     `json:"slot_id"`
	FeedBackId    int     `json:"feedback_id,omitempty"`
	ListServiceId []int   `json:"list_service"`
}
type BookingRequest struct {
	CustomerID    int     `json:"customer_id"`
	BarberId      int     `json:"barber_id"`
	ResultId      int     `json:"result_id,omitempty"`
	Status        string  `json:"status"`
	Price         float32 `json:"price"`
	SlotId        int     `json:"slot_id"`
	FeedBackId    int     `json:"feedback_id,omitempty"`
	ListServiceId []int   `json:"list_service"`
}
