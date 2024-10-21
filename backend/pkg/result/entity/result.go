package entity

type Result struct {
	Id          int    `json:"id"`
	BarberId    int    `json:"barber_id"`
	UserId      int    `json:"user_id"`
	BookingId   int    `json:"booking_id"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type ImageResult struct {
	Id        int    `json:"id"`
	Url       string `json:"url"`
	ResultId  int    `json:"result_id"`
	CreatedAt int64  `json:"created_at"`
}
