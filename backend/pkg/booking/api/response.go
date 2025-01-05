package api

type BookingResponse struct {
	ID           int               `json:"id"`
	CustomerID   int               `json:"customer_id"`
	CustomerName string            `json:"customer_name"`
	BarberId     int               `json:"barber_id"`
	BarberName   string            `json:"barber_name"`
	ResultId     int               `json:"result_id"`
	Status       string            `json:"status"`
	Price        float32           `json:"price"`
	SlotId       int               `json:"slot_id"`
	BookedDate   string            `json:"booked_date"`
	StartTime    string            `json:"start_time"`
	FeedBackId   int               `json:"feedback_id"`
	CreatedAt    int64             `json:"created_at"`
	UpdatedAt    int64             `json:"updated_at"`
	ListServices []ServiceResponse `json:"list_services"`
	PreviewId    int               `json:"preview_id,omitempty"`
}
type ServiceResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Url         string `json:"url"`
}
type KafkaBookingResponse struct {
	UUID         string  `json:"uuid"`
	ID           int     `json:"id"`
	CustomerID   int     `json:"customer_id"`
	BarberId     int     `json:"barber_id"`
	ResultId     int     `json:"result_id"`
	Status       string  `json:"status"`
	Price        float32 `json:"price"`
	SlotId       int     `json:"slot_id"`
	FeedBackId   int     `json:"feedback_id"`
	CreatedAt    int64   `json:"created_at"`
	UpdatedAt    int64   `json:"updated_at"`
	ListServices []int32 `json:"list_services"`
}
