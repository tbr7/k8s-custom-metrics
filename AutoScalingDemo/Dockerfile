FROM golang:1.22-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s' -o metrics-exporter

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/metrics-exporter .
COPY cert.pem .
COPY key.pem .

EXPOSE 5001

CMD ["./metrics-exporter"]
