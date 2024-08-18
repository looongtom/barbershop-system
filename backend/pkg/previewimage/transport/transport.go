package transport

import (
	"DoAn/pkg/previewimage"
	"DoAn/pkg/previewimage/api"
	"DoAn/pkg/previewimage/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"net/http"
)

type (
	GetPreviewImageRequest struct {
		ID int `json:"id"`
	}

	Response struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func MakeCreatePreviewImageEndpoints(svc previewimage.PreviewImageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.CreatePreviewImageRequest)
		resp, err := svc.CreatePreviewImage(ctx, api.CreatePreviewImageRequest{
			Url:       req.Url,
			AccountId: req.AccountId})
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeGetPreviewImageEndpoints(svc service.PreviewImageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetPreviewImageRequest)
		resp, err := svc.GetPreviewImage(ctx, req.ID)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeGetListPreviewImageByAccountIdEndpoints(svc service.PreviewImageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetPreviewImageRequest)
		resp, err := svc.GetListPreviewImageByAccountId(ctx, req.ID)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func DecodeCreatePreviewImageRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.CreatePreviewImageRequest
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

	req.Url = file
	return req, err
}

func DecodeGetPreviewImageRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req GetPreviewImageRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
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
