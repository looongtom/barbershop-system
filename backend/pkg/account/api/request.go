package api

type CheckExistedBarberRequest struct {
	UserId int `json:"user_id"`
}
type VerifyTokenRequest struct {
	Token string `json:"token"`
}
