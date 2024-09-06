package entity

type Criteria struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Img        string `json:"img"`
	CategoryId int    `json:"category_id"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

type CategoryCriteria struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
