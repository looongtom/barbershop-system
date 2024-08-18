package entity

type PreviewImage struct {
	ID        int    `json:"id,omitempty"`
	Url       string `json:"url"`
	CreatedAt int64  `json:"created_at,omitempty"`
	AccountId int    `json:"account_id"`
}
