# docker build -t simple-app .
# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.22 as builder

WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.mod .
COPY go.sum* ./

RUN go mod download

# Copy local code to the container image.
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s' -o myapp

# Alpine for its small footprint compared to Ubuntu.
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/myapp .

EXPOSE 8080

# Run the web service on container startup.
CMD ["./myapp"]
