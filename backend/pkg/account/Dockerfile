FROM golang:1.23-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o ./cmd/apiserver ./cmd/main.go


FROM scratch

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder /app/cmd/apiserver /
COPY /cmd/account.env .

# Command to run when starting the container.
ENTRYPOINT ["/apiserver"]