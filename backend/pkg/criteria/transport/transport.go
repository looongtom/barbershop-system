package transport

import (
	"DoAn/pkg/criteria"
	"DoAn/pkg/criteria/api"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
)

type (
	CreateCategoryRequest struct {
		Name string `json:"name"`
	}
	CreateOrUpdateCriteriaRequest struct {
		Id         int    `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
		Img        string `json:"img,omitempty"`
		CategoryId int    `json:"category_id,omitempty"`
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
		req := request.(api.CreateOrUpdateCategory)
		resp, err := svc.CreateOrUpdateCategory(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeCreateOrUpdateCriteriaEndpoints(svc criteria.CriteriaService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.CreateOrUpdateCriteria)
		resp, err := svc.CreateOrUpdateCriteria(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeGetCriteriaEndpoints(svc criteria.CriteriaService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		resp, err := svc.GetCriteria(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeGetCategoryEndpoints(svc criteria.CriteriaService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		resp, err := svc.GetCategory(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeGetListCriteriaEndpoints(svc criteria.CriteriaService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp, err := svc.GetListCriteria(ctx)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeGetListCategoryEndpoints(svc criteria.CriteriaService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp, err := svc.GetListCategory(ctx)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeDeleteCriteriaEndpoints(svc criteria.CriteriaService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		resp, err := svc.DeleteCriteria(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}

}

func MakeDeleteCategoryEndpoints(svc criteria.CriteriaService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		resp, err := svc.DeleteCategory(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func DecodeEmptyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func DecodeCreateOrUpdateCategoryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request api.CreateOrUpdateCategory
	fmt.Println(request)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeCreateOrUpdateCriteriaRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request api.CreateOrUpdateCriteria
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		return nil, err
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	request.Img = file
	categoryIdValue, ok := strconv.Atoi(r.FormValue("category_id"))
	if ok != nil {
		fmt.Println("Error Retrieving category_id")
		return nil, ok
	}
	idValue := r.FormValue("id")
	if idValue != "" && &idValue != nil {
		id, ok := strconv.Atoi(idValue)
		if ok != nil {
			fmt.Println("Error Retrieving id")
			return nil, ok
		}
		request.ID = id
	}
	request.Name = r.FormValue("name")
	request.CategoryId = categoryIdValue
	return request, nil
}

func DecodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := r.URL.Query().Get("id")
	return id, nil
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
