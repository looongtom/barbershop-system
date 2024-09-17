package transport

import (
	"DoAn/pkg/booking"
	"DoAn/pkg/booking/api"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type (
	CreateBookingRequest struct {
		BarberId   string  `json:"barber_id"`
		CustomerID string  `json:"customer_id"`
		FeedBackId string  `json:"feedback_id,omitempty"`
		SlotId     string  `json:"slot_id"`
		ResultId   string  `json:"result_id,omitempty"`
		Status     string  `json:"status"`
		Price      float32 `json:"price"`
	}
	UpdateBookingRequest struct {
		Id         string  `json:"id"`
		BarberId   string  `json:"barber_id"`
		CustomerID string  `json:"customer_id"`
		FeedBackId string  `json:"feedback_id,omitempty"`
		SlotId     string  `json:"slot_id"`
		ResultId   string  `json:"result_id,omitempty"`
		Status     string  `json:"status"`
		Price      float32 `json:"price"`
	}

	Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func MakeGetBookingEndpoints(svc booking.BookingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		resp, err := svc.GetBooking(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeGetListBookingEndpoints(svc booking.BookingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp, err := svc.GetListBooking(ctx)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeCreateBookingEndpoints(svc booking.BookingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.BookingRequest)
		resp, err := svc.CreateBooking(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}
func MakeCreateBookingKafkaEndpoints(svc booking.BookingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.BookingRequest)
		resp, err := svc.CreateBookingKafka(ctx, req)
		if resp == nil {
			return Response{
				Message: "failed to create booking",
				Data:    nil,
			}, nil
		} else if resp != nil && resp.(api.KafkaBookingResponse).ID == 0 {
			return Response{
				Message: "failed to create booking",
				Data:    nil,
			}, nil
		}
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeUpdateBookingEndpoints(svc booking.BookingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.UpdateBookingRequest)
		resp, err := svc.UpdateBooking(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func DecodeCreateBookingRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req api.BookingRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeGetBookingRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	id := r.URL.Query().Get("id")
	return id, nil
}

func DecodeEmptyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
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

func DecodeUpdateBookingRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
	var req api.UpdateBookingRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
