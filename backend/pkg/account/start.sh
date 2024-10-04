#!/bin/sh

# Start the HTTP server
./main &

# Start the gRPC server
./grpc_server &

# Wait for all background processes to finish
wait