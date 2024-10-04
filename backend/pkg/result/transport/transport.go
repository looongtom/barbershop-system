package transport

import (
	"DoAn"
	"DoAn/api"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type (
	CreateOrUpdateResultRequest struct {
		ID          int    `json:"id,omitempty"`
		BarberId    int    `json:"barber_id,omitempty"`
		UserId      int    `json:"user_id,omitempty"`
		Description string `json:"description,omitempty"`
	}

	Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func MakeCreateOrUpdateResultEndpoints(svc result.ResultService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.CreateOrUpdateResult)
		resp, err := svc.CreateOrUpdateResult(ctx, req)
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

func MakeGetResultByBarberIdEndpoints(svc result.ResultService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.GetResultById)
		resp, err := svc.GetResultByBarberId(ctx, req.ID)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}
func MakeGetResultByUserIdEndpoints(svc result.ResultService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.GetResultById)
		resp, err := svc.GetResultByBarberId(ctx, req.ID)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}
func MakeGetResultByBookingIdEndpoints(svc result.ResultService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.GetResultById)
		resp, err := svc.GetResultByBarberId(ctx, req.ID)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func DecodeCreateOrUpdateResult(_ context.Context, r *http.Request) (interface{}, error) {
	var request api.CreateOrUpdateResult
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	if &request.ID == nil {
		if request.Description == "" || &request.Description == nil {
			return nil, errors.New("description is required")
		}
	}
	return request, nil
}

func DecodeGetResultById(_ context.Context, r *http.Request) (interface{}, error) {
	var request api.GetResultById
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	if &request.ID == nil {
		return nil, errors.New("id is required")
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
