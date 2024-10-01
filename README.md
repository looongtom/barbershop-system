"# barbershop-system" 



go run ./pkg/account/grpc/grpc_server.go

go run ./pkg/timeslot/grpc/grpc_server.go

go run ./pkg/servicing/grpc/grpc_server.go

go run ./pkg/booking/grpc/grpc_server.go

go run ./pkg/account/cmd/main.go

go run ./pkg/booking/cmd/main.go

go run ./pkg/criteria/cmd/main.go

go run ./pkg/previewimage/cmd/main.go

go run ./pkg/servicing/cmd/main.go

go run ./pkg/timeslot/cmd/main.go

go run ./pkg/notification/previewimage-listener/consumer_previewimage.go

go run ./pkg/result_preview_ws/server_ws.go

go run ./pkg/notification/hairfast-listener/hairfast.go

