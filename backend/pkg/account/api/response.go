package api

type GetAccountResponse struct {
	ID          int    `json:"id,omitempty"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Role        int    `json:"role,omitempty"`
	PhoneNumber string `json:"phoneNumber"`
	FullName    string `json:"fullName"`
	About       string `json:"about"`
	Dob         string `json:"dob"`
	Avatar      string `json:"avatar"`
	CreatedAt   int64  `json:"created_at,omitempty"`
	UpdatedAt   int64  `json:"updated_at,omitempty"`
}
