package api

import (
	"mime/multipart"
)

type CreatePreviewImageRequest struct {
	Url       multipart.File `json:"url"`
	AccountId int            `json:"account_id"`
}
