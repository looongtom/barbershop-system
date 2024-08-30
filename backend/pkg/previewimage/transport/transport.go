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
	"strconv"
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
		resp, err := svc.CreatePreviewImage(ctx, req)
		return Response{
			Message: "success",
			Data:    resp,
		}, err
	}
}

func MakeUploadImagesEndpoints(svc previewimage.PreviewImageService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.UpdateImageRequest)
		resp, err := svc.UploadImages(ctx, req)
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

func DecodeUploadImagesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req api.UpdateImageRequest
	r.ParseMultipartForm(10 << 20)

	selfImg, handler1, err := r.FormFile("self_img")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		return nil, err
	}
	defer selfImg.Close()
	fmt.Printf("Uploaded File: %+v\n", handler1.Filename)
	fmt.Printf("File Size: %+v\n", handler1.Size)
	fmt.Printf("MIME Header: %+v\n", handler1.Header)

	shapeImg, handler2, err := r.FormFile("shape_img")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		return nil, err
	}
	defer selfImg.Close()
	fmt.Printf("Uploaded File: %+v\n", handler2.Filename)
	fmt.Printf("File Size: %+v\n", handler2.Size)
	fmt.Printf("MIME Header: %+v\n", handler2.Header)

	colorImg, handler3, err := r.FormFile("color_img")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		return nil, err
	}
	defer selfImg.Close()
	fmt.Printf("Uploaded File: %+v\n", handler3.Filename)
	fmt.Printf("File Size: %+v\n", handler3.Size)
	fmt.Printf("MIME Header: %+v\n", handler3.Header)

	req.SelfImg = selfImg
	req.ShapeImg = shapeImg
	req.ColorImg = colorImg

	accountStr := r.FormValue("account_id")
	accountValue, ok := strconv.Atoi(accountStr)
	if ok != nil {
		fmt.Println("Error Retrieving account_id")
		return nil, ok
	}
	req.AccountId = accountValue
	return req, err
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
	accountStr := r.FormValue("account_id")
	accountValue, ok := strconv.Atoi(accountStr)
	if ok != nil {
		fmt.Println("Error Retrieving account_id")
		return nil, ok
	}
	req.AccountId = accountValue
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
