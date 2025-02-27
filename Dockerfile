FROM golang:1.22.8-alpine AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download
RUN go mod verify

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o charts-proxy

FROM alpine:3

COPY --from=builder /app/charts-proxy /usr/local/bin/charts-proxy

ENTRYPOINT ["/usr/local/bin/charts-proxy"]
