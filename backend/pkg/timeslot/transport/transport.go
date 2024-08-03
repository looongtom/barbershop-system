package transport

import (
	"DoAn/pkg/timeslot"
	"DoAn/pkg/timeslot/api"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"net/http"
)

type (
	CreateTimeSlotRequest struct {
		BarberId   int    `json:"barber_id"`
		StartTime  string `json:"start_time"`
		BookedDate string `json:"booked_date"`
		Status     string `json:"status"`
	}

	Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func MakeGetTimeSlotEndpoints(svc timeslot.TimeSlotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.FindTimeslotRequest)
		resp, err := svc.GetListTimeSlotByBarberId(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeCreateTimeSlotEndpoints(svc timeslot.TimeSlotService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.CreateOrUpdateTimeslotRequest)
		resp, err := svc.CreateTimeSlot(ctx, req)
		message := "updated successfully"
		if req.ID == 0 {
			message = "created successfully"
		}
		return Response{
			Message: message,
			Data:    resp,
		}, err
	}
}

func DecodeCreateTimeSlotRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request api.CreateOrUpdateTimeslotRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	if &request.ID == nil {
		if request.StartTime == "" || &request.StartTime == nil {
			return nil, errors.New("missing start time")
		}
		if request.BookedDate == "" || &request.BookedDate == nil {
			return nil, errors.New("missing BookedDate")
		}
		if request.Status == "" || &request.Status == nil {
			return nil, errors.New("missing Status")
		}
		if &request.BarberId == nil {
			return nil, errors.New("missing BarberId")
		}
	}
	return request, nil
}

func DecodeUpdatedTimeSlotRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request api.CreateOrUpdateTimeslotRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeGetTimeSlotRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request api.FindTimeslotRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err, ok := response.(error); ok {
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(response)
}
