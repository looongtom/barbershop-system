package api

import "mime/multipart"

type CreateOrUpdateCriteria struct {
	ID         int            `json:"id,omitempty"`
	Name       string         `json:"name,omitempty"`
	Img        multipart.File `json:"img,omitempty"`
	CategoryId int            `json:"category_id,omitempty"`
}

type CreateOrUpdateCategory struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type FindCriteria struct {
	Name       string `json:"name,omitempty"`
	CategoryId int    `json:"category_id,omitempty"`
}
