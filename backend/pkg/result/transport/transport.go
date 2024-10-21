package transport

import (
	"DoAn"
	"DoAn/api"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
		resp, err := svc.CreateResult(ctx, req)
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
		req := request.(int)
		resp, err := svc.GetResultByBookingId(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func DecodeCreateOrUpdateResult(_ context.Context, r *http.Request) (interface{}, error) {
	var request api.CreateOrUpdateResult
	r.ParseMultipartForm(10 << 20)
	//get list image
	files := r.MultipartForm.File["list_img"]
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			return nil, err
		}
		request.ListImg = append(request.ListImg, f)
	}
	barberIdValue, ok := strconv.Atoi(r.FormValue("barber_id"))
	if ok != nil {
		fmt.Println("Error Retrieving category_id")
		return nil, ok
	}
	cusIdValue, ok := strconv.Atoi(r.FormValue("user_id"))
	if ok != nil {
		fmt.Println("Error Retrieving category_id")
		return nil, ok
	}
	bookingIdValue, ok := strconv.Atoi(r.FormValue("booking_id"))
	if ok != nil {
		fmt.Println("Error Retrieving category_id")
		return nil, ok
	}
	request.Description = r.FormValue("description")

	idValue := r.FormValue("id")
	if idValue != "" && &idValue != nil {
		id, ok := strconv.Atoi(idValue)
		if ok != nil {
			fmt.Println("Error Retrieving id")
			return nil, ok
		}
		request.ID = id
	}
	request.BarberId = barberIdValue
	request.UserId = cusIdValue
	request.BookingId = bookingIdValue
	return request, nil
}

func DecodeGetResultById(_ context.Context, r *http.Request) (interface{}, error) {
	id := r.URL.Query().Get("id")
	if id != "" && &id != nil {
		idValue, ok := strconv.Atoi(id)
		if ok != nil {
			fmt.Println("Error Retrieving id")
			return nil, ok
		}
		return idValue, nil
	}
	return nil, errors.New("id is required")
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
