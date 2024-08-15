"# barbershop-system" 



go run ./pkg/account/grpc/grpc_server.go

go run ./pkg/timeslot/grpc/grpc_server.go

go run ./pkg/servicing/grpc/grpc_server.go

go run ./pkg/booking/grpc/grpc_server.go