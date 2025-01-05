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
	BookedDate    string  `json:"booked_date"`
	PreviewId     int     `json:"preview_id,omitempty"`
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

type GetBookingByIdRequest struct {
	Id int `json:"id"`
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

type UpdateBookingStatusRequest struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}

type UpdateBookingTimeslotRequest struct {
	Id         int `json:"id"`
	TimeslotId int `json:"timeslot_id"`
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
