package api

type TimeslotRequest struct {
	BarberId   int    `json:"barber_id"`
	StartTime  string `json:"start_time"`
	BookedDate string `json:"booked_date"`
	Status     string `json:"status"`
}
type FindTimeslotRequest struct {
	BarberId   int    `json:"barber_id"`
	StartTime  string `json:"start_time"`
	BookedDate string `json:"booked_date"`
}

type CheckExistTimeslotRequest struct {
	Id int `json:"id"`
}

type UpdateTimeslotRequest struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type CreateOrUpdateTimeslotRequest struct {
	ID         int    `json:"id,omitempty"`
	BarberId   int    `json:"barber_id,omitempty"`
	StartTime  string `json:"start_time,omitempty"`
	BookedDate string `json:"booked_date,omitempty"`
	Status     string `json:"status,omitempty"`
}

type CreateListTimeslotRequest struct {
	ListTimeslot []CreateOrUpdateTimeslotRequest `json:"list"`
}

type CheckExistedTimeslotRequest struct {
	BarberId   int    `json:"barber_id,omitempty"`
	StartTime  string `json:"start_time,omitempty"`
	BookedDate string `json:"booked_date,omitempty"`
}
