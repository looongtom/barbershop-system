package api

import "mime/multipart"

type GetDetailResultResponse struct {
	ID          int                 `json:"id,omitempty"`
	BarberId    int                 `json:"barber_id,omitempty"`
	UserId      int                 `json:"user_id,omitempty"`
	BookingId   int                 `json:"booking_id,omitempty"`
	Description string              `json:"description,omitempty"`
	ListImg     []GetResultResponse `json:"list_img,omitempty"`
}

type GetResultResponse struct {
	ID       int    `json:"id,omitempty"`
	Url      string `json:"url,omitempty"`
	ResultId int    `json:"result_id,omitempty"`
}

type CreateOrUpdateResultResponse struct {
	ID          int                           `json:"id,omitempty"`
	BarberId    int                           `json:"barber_id,omitempty"`
	UserId      int                           `json:"user_id,omitempty"`
	BookingId   int                           `json:"booking_id,omitempty"`
	Description string                        `json:"description,omitempty"`
	ListImg     []CreateOrUpdateImageResponse `json:"list_img,omitempty"`
}

type CreateOrUpdateImageResponse struct {
	ID       int    `json:"id,omitempty"`
	Url      string `json:"url,omitempty"`
	ResultId int    `json:"result_id,omitempty"`
}

type CreateOrUpdateResult struct {
	ID          int              `json:"id,omitempty"`
	BarberId    int              `json:"barber_id,omitempty"`
	UserId      int              `json:"user_id,omitempty"`
	BookingId   int              `json:"booking_id,omitempty"`
	Description string           `json:"description,omitempty"`
	ListImg     []multipart.File `json:"list_img,omitempty"`
}

type CreateImageResult struct {
	Url      string `json:"url,omitempty"`
	ResultId int    `json:"result_id,omitempty"`
}

type GetResultById struct {
	ID int `json:"id"`
}
