package mapper

type BookingMapper struct {
	ID                int              `json:"id"`
	CustomerID        int              `json:"customer_id"`
	CustomerName      string           `json:"customer_name"`
	BarberId          int              `json:"barber_id"`
	BarberName        string           `json:"barber_name"`
	ResultId          int              `json:"result_id"`
	Status            string           `json:"status"`
	Price             float32          `json:"price"`
	SlotId            int              `json:"slot_id"`
	TimeSlot          BookingTimeSlot  `json:"time_slot"`
	FeedBackId        int              `json:"feedback_id"`
	CreatedAt         int64            `json:"created_at"`
	UpdatedAt         int64            `json:"updated_at"`
	ListServices      []int64          `json:"list_services"`
	ListServiceStruct []BookingService `json:"list_service_struct"`
	BookedDate        string           `json:"booked_date"`
}

type BookingService struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

type BookingTimeSlot struct {
	ID         int    `json:"id"`
	StartTime  string `json:"start_time"`
	BookedDate string `json:"booked_date"`
}
