package api

type BookingRequest struct {
	CustomerID    int     `json:"customer_id"`
	BarberId      int     `json:"barber_id"`
	ResultId      int     `json:"result_id,omitempty"`
	Status        string  `json:"status"`
	Price         float32 `json:"price"`
	SlotId        int     `json:"slot_id"`
	FeedBackId    int     `json:"feedback_id,omitempty"`
	ListServiceId []int   `json:"list_service"`
}

type UpdateBookingRequest struct {
	Id         int     `json:"id"`
	CustomerID int     `json:"customer_id"`
	BarberId   int     `json:"barber_id"`
	ResultId   int     `json:"result_id,omitempty"`
	Status     string  `json:"status"`
	Price      float32 `json:"price"`
	SlotId     int     `json:"slot_id"`
	FeedBackId int     `json:"feedback_id,omitempty"`
}

type FindBookingRequest struct {
	CustomerId int `json:"customerId"`
	BarberId   int `json:"barberId"`
}
