package api

type GetServiceByIdResponse struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Description  string `json:"description"`
	Url          string `json:"url"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
}
