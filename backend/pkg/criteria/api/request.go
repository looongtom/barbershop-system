package api

type CreateOrUpdateCriteria struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Img        string `json:"img,omitempty"`
	CategoryId int    `json:"category_id,omitempty"`
}

type CreateOrUpdateCategory struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
