package entity

type Token struct {
	Token     string `json:"token"`
	User      string `json:"account"`
	CreatedAt string `json:"created_at"`
	ExpiredAt string `json:"expired_at"`
}
