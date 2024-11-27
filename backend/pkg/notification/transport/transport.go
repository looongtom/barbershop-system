package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	notification "DoAn"
	"DoAn/api"
)

type getListNotificationResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func MakeGetListNotificationEndpoints(svc notification.NotificationService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.GetListNotificationRequest)
		resp, err := svc.GetNotification(ctx, req)
		if err != nil {
			return getListNotificationResponse{nil, err.Error()}, nil
		}
		return getListNotificationResponse{resp, "success"}, nil
	}
}

func DecodeGetListNotificationRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req api.GetListNotificationRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
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
