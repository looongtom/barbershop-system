package entity

type PreviewImage struct {
	ID           int    `json:"id,omitempty"`
	GeneratedImg string `json:"generated_img,omitempty"`
	CreatedAt    int64  `json:"created_at,omitempty"`
	AccountId    int    `json:"account_id"`
	SelfImg      string `json:"self_img"`
	ShapeImg     string `json:"shape_img"`
	ColorImg     string `json:"color_img"`
}
