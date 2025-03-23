# Build the app
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o /app/app ./cmd/api

# Execute the app
FROM alpine

WORKDIR /app
COPY --from=builder /app/app .

CMD ["/app/app"]
