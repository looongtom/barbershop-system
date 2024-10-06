package mapper

type BookingMapper struct {
	ID                int              `json:"id"`
	CustomerID        int              `json:"customer_id"`
	BarberId          int              `json:"barber_id"`
	ResultId          int              `json:"result_id"`
	Status            string           `json:"status"`
	Price             float32          `json:"price"`
	SlotId            int              `json:"slot_id"`
	FeedBackId        int              `json:"feedback_id"`
	CreatedAt         int64            `json:"created_at"`
	UpdatedAt         int64            `json:"updated_at"`
	ListServices      []int64          `json:"list_services"`
	ListServiceStruct []BookingService `json:"list_service_struct"`
}

type BookingService struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Url         string `json:"url"`
}
