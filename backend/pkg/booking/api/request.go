package api

type VerifyTokenRequest struct {
	Token string `json:"token"`
}
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
type KafkaBookingRequest struct {
	UUID          string  `json:"uuid"`
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
	Id            int     `json:"id"`
	CustomerID    int     `json:"customer_id"`
	BarberId      int     `json:"barber_id"`
	ResultId      int     `json:"result_id,omitempty"`
	Status        string  `json:"status"`
	Price         float32 `json:"price"`
	SlotId        int     `json:"slot_id"`
	FeedBackId    int     `json:"feedback_id,omitempty"`
	ListServiceId []int   `json:"list_service"`
}

type UpdateBookingServiceRequest struct {
	Id            int   `json:"id"`
	ListServiceId []int `json:"list_service"`
}

type FindBookingRequest struct {
	CustomerId int `json:"customerId"`
	BarberId   int `json:"barberId"`
}

type FindListBookingRequest struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"pageSize"`
	Account    int    `json:"account_id"`
	BookedDate string `json:"booked_date"`
	Role       string `json:"role"`
}
