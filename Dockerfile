# Build stage
FROM golang:1.22 AS builder
WORKDIR /app

# Cache and install dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY Makefile .
COPY doc.go .
COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal

# Build the application
RUN make build

# Final stage
FROM alpine:latest
WORKDIR /app/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/statelessdb .

EXPOSE 8080

# Command to run
CMD ["./statelessdb"]
