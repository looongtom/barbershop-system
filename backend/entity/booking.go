package entity

type Booking struct {
	ID         int     `json:"id"`
	CustomerID int     `json:"customer_id"`
	BarberId   int     `json:"barber_id"`
	ResultId   int     `json:"result_id"`
	Status     string  `json:"status"`
	Price      float32 `json:"price"`
	SlotId     int     `json:"slot_id"`
	FeedBackId int     `json:"feedback_id"`
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
