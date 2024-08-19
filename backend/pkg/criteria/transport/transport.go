package transport

import (
	"DoAn/pkg/criteria"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"net/http"
)

type (
	CreateCategoryRequest struct {
		Name string `json:"name"`
	}
	CreateCriteriaRequest struct {
		Name       string `json:"name"`
		Img        string `json:"img"`
		CategoryId int    `json:"category_id"`
	}
	CreateCategory struct {
		Name string `json:"name"`
	}
	UpdateCriteriaRequest struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Img        string `json:"img"`
		CategoryId int    `json:"category_id"`
	}
	Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func MakeCreateCategoryEndpoints(svc criteria.CriteriaService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		resp, err := svc.CreateCategory(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func DecodeEmptyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func DecodeCreateCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request CreateCategoryRequest
	fmt.Println(request)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
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
