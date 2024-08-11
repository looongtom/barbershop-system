package api

type BookingResponse struct {
	ID           int      `json:"id"`
	CustomerID   int      `json:"customer_id"`
	BarberId     int      `json:"barber_id"`
	ResultId     int      `json:"result_id"`
	Status       string   `json:"status"`
	Price        float32  `json:"price"`
	SlotId       int      `json:"slot_id"`
	FeedBackId   int      `json:"feedback_id"`
	CreatedAt    int64    `json:"created_at"`
	UpdatedAt    int64    `json:"updated_at"`
	ListServices []string `json:"list_services"`
}
