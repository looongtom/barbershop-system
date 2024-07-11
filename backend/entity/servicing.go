package entity

type Servicing struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Url         string `json:"url"`
	CategoryID  int    `json:"category_id"`

	CreatedAt int64 `json:"created_at,omitempty"`
	UpdatedAt int64 `json:"updated_at,omitempty"`
}

type Category struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`

	CreatedAt int64 `json:"created_at,omitempty"`
	UpdatedAt int64 `json:"updated_at,omitempty"`
}
