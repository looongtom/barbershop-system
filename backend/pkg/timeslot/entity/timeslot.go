package entity

type Timeslot struct {
	ID         int    `json:"id"`
	StartTime  string `json:"start_time"`
	BookedDate string `json:"booked_date"`
	Status     string `json:"status"`
	BarberId   int    `json:"barber_id"`
	CreatedAt  int64  `json:"created_at,omitempty"`
	UpdatedAt  int64  `json:"updated_at,omitempty"`
}
