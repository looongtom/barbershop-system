
*generate account*

protoc --go_out=. --go_opt=paths=source_relative account.proto
protoc --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative account.proto