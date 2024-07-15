package entity

type Booking struct {
	ID         string  `json:"id"`
	CustomerID string  `json:"customer_id"`
	BarberId   string  `json:"barber_id"`
	ResultId   string  `json:"result_id"`
	Status     string  `json:"status"`
	Price      float32 `json:"price"`
	SlotId     string  `json:"slot_id"`
	FeedBackId string  `json:"feedback_id"`
	CreatedAt  int64   `json:"created_at"`
	UpdatedAt  int64   `json:"updated_at"`
}

type BookingDetail struct {
	ID        string `json:"id"`
	BookingId string `json:"booking_id"`
	ServiceId string `json:"service_id"`
}

type TimeSlot struct {
	ID         string `json:"id"`
	StartTime  int64  `json:"start_time"`
	BookedDate int64  `json:"booked_date"`
	Status     string `json:"status"`
	BarberId   string `json:"barber_id"`
}
