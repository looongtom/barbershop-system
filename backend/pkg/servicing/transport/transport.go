package transport

import (
	"DoAn"
	"DoAn/entity"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type (
	CreateServiceRequest struct {
		Service entity.Servicing
	}
	CreateCategoryRequest struct {
		Name string `json:"name"`
	}
	UpdateCategoryRequest struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func MakeGetServiceEndpoints(svc servicing.ServicingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		resp, err := svc.GetServicing(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeGetCategoryEndpoints(svc servicing.ServicingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		resp, err := svc.GetCategory(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeGetListCategoryEndpoints(svc servicing.ServicingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp, err := svc.GetListCategory(ctx)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}
func MakeGetListServiceEndpoints(svc servicing.ServicingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp, err := svc.GetListServicing(ctx)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeGetListServiceAndCategoryEndpoints(svc servicing.ServicingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp, err := svc.GetListServicingAndCategory(ctx)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeCreateServiceEndpoints(svc servicing.ServicingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateServiceRequest)
		resp, err := svc.CreateServicing(ctx, req.Service)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}
func MakeUpdateServiceEndpoints(svc servicing.ServicingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateServiceRequest)
		resp, err := svc.UpdateServicing(ctx, req.Service)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}
func MakeUpdateCategoryEndpoints(svc servicing.ServicingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateCategoryRequest)
		categoryReq := entity.Category{
			ID:   req.ID,
			Name: req.Name,
		}
		resp, err := svc.UpdateCategory(ctx, categoryReq)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeCreateCategoryEndpoints(svc servicing.ServicingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateCategoryRequest)
		fmt.Println(req)
		resp, err := svc.CreateCategory(ctx, req.Name)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func DecodeCreateCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request CreateCategoryRequest
	fmt.Println(request)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
func DecodeUpdateCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request UpdateCategoryRequest
	fmt.Println(request)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeCreateServiceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&request.Service); err != nil {
		return nil, err
	}
	return request, nil
}
func DecodeGetServiceRequest(_ context.Context, r *http.Request) (interface{}, error) {
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
