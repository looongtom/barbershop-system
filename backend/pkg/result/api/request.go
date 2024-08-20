package api

type CreateOrUpdateResult struct {
	ID          int    `json:"id,omitempty"`
	BarberId    int    `json:"barber_id,omitempty"`
	UserId      int    `json:"user_id,omitempty"`
	BookingId   int    `json:"booking_id,omitempty"`
	Description string `json:"description,omitempty"`
}

type CreateImageResult struct {
	ID          int    `json:"id,omitempty"`
	Url         string `json:"url,omitempty"`
	ResultId    int    `json:"result_id,omitempty"`
	Description string `json:"description,omitempty"`
}

type GetResultById struct {
	ID int `json:"id"`
}
