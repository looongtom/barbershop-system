package api

import (
	"mime/multipart"
)

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
